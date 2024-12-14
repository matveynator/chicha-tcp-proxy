package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

// Route describes a single forwarding route configuration.
type Route struct {
	LocalPort  string
	RemoteIP   string
	RemotePort string
}

func main() {
	routesFlag := flag.String("routes", "", "Comma-separated list of routes in the format LOCALPORT:REMOTEIP:REMOTEPORT")
	logFile := flag.String("log", "chicha-tcp-proxy.log", "Path to the log file")
	rotationFrequency := flag.Duration("rotation", 24*time.Hour, "Log rotation frequency (e.g. 24h, 1h, etc.)")
	flag.Parse()

	if *routesFlag == "" {
		log.Fatal("Error: The -routes flag is required.")
	}

	routes, err := parseRoutes(*routesFlag)
	if err != nil {
		log.Fatalf("Error parsing routes: %v", err)
	}
	if len(routes) == 0 {
		log.Fatalf("Error: no valid routes found in '%s'", *routesFlag)
	}

	// Print routes, log file path, and rotation frequency
	fmt.Println("========== CHICHA TCP PROXY ==========")
	fmt.Println("Routes:")
	for _, route := range routes {
		fmt.Printf("  LocalPort=%s -> RemoteIP=%s RemotePort=%s\n", route.LocalPort, route.RemoteIP, route.RemotePort)
	}
	fmt.Printf("Log file: %s\n", *logFile)
	fmt.Printf("Log rotation frequency: %v\n", *rotationFrequency)
	fmt.Println("======================================")

	logger, file, err := setupLogger(*logFile)
	if err != nil {
		log.Fatalf("Error setting up logger: %v", err)
	}

	log.Printf("Starting chicha-tcp-proxy")

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	logger.Printf("Using %d CPU cores", numCPUs)
	log.Printf("Using %d CPU cores", numCPUs)

	// Start log rotation in a separate goroutine.
	go rotateLogs(*logFile, file, logger, *rotationFrequency)

	// Start proxy servers for each route.
	for _, route := range routes {
		logger.Printf("Starting proxy for route: local=%s remote=%s:%s", route.LocalPort, route.RemoteIP, route.RemotePort)
		log.Printf("Starting proxy for route: local=%s remote=%s:%s", route.LocalPort, route.RemoteIP, route.RemotePort)
		go startProxy(":"+route.LocalPort, route.RemoteIP+":"+route.RemotePort, logger)
	}

	// Block forever.
	select {}
}

// parseRoutes parses the routes from a string like "8080:46.4.70.114:80,8443:46.4.70.114:443"
func parseRoutes(routesFlag string) ([]Route, error) {
	parts := strings.Split(routesFlag, ",")
	var routes []Route
	for _, part := range parts {
		segments := strings.Split(part, ":")
		if len(segments) != 3 {
			return nil, fmt.Errorf("invalid route format: '%s' (expected LOCALPORT:REMOTEIP:REMOTEPORT)", part)
		}
		routes = append(routes, Route{
			LocalPort:  segments[0],
			RemoteIP:   segments[1],
			RemotePort: segments[2],
		})
	}
	return routes, nil
}

// setupLogger creates or opens the log file and returns a logger and the file handle.
func setupLogger(logFile string) (*log.Logger, *os.File, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file '%s': %v", logFile, err)
	}
	logger := log.New(file, "", log.LstdFlags)
	return logger, file, nil
}

// rotateLogs handles log rotation based on the specified frequency.
func rotateLogs(logFile string, file *os.File, logger *log.Logger, frequency time.Duration) {
	for {
		time.Sleep(frequency)

		file.Close()

		rotatedFile := logFile + "." + time.Now().Format("2006-01-02")
		if err := os.Rename(logFile, rotatedFile); err != nil {
			logger.Printf("Error rotating logs: %v", err)
			log.Printf("Error rotating logs: %v", err) // Duplicate log to console
			// Attempt to reopen the current log file
			newFile, err2 := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err2 != nil {
				logger.Fatalf("Failed to reopen log file after rotation error: %v", err2)
				log.Fatalf("Failed to reopen log file after rotation error: %v", err2) // Duplicate log to console
			}
			file = newFile
			logger.SetOutput(file)
			continue
		}

		// Open a new log file for future logging
		newFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatalf("Failed to create new log file after rotation: %v", err)
			log.Fatalf("Failed to create new log file after rotation: %v", err) // Duplicate log to console
		}
		file = newFile
		logger.SetOutput(file)
		logger.Println("Log file rotated successfully, now compressing old log...")
		log.Println("Log file rotated successfully, now compressing old log...") // Duplicate log to console

		// Compress the rotated file
		if err := compressFile(rotatedFile); err != nil {
			logger.Printf("Error compressing rotated file: %v", err)
			log.Printf("Error compressing rotated file: %v", err) // Duplicate log to console
		} else {
			logger.Printf("Compression successful: %s.gz", rotatedFile)
			log.Printf("Compression successful: %s.gz", rotatedFile) // Duplicate log to console
			if err := os.Remove(rotatedFile); err != nil {
				logger.Printf("Error removing uncompressed rotated file: %v", err)
				log.Printf("Error removing uncompressed rotated file: %v", err) // Duplicate log to console
			}
		}
	}
}

// compressFile compresses the specified file into a .gz file.
func compressFile(filename string) error {
	original, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for compression: %v", err)
	}
	defer original.Close()

	gzFile, err := os.OpenFile(filename+".gz", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create gz file: %v", err)
	}
	defer gzFile.Close()

	gzWriter := gzip.NewWriter(gzFile)
	defer gzWriter.Close()

	if _, err := io.Copy(gzWriter, original); err != nil {
		return fmt.Errorf("failed to copy data for compression: %v", err)
	}

	return nil
}

// startProxy starts a TCP listener on listenAddr and forwards traffic to targetAddr.
func startProxy(listenAddr, targetAddr string, logger *log.Logger) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Failed to start proxy on %s: %v", listenAddr, err)
	}
	defer listener.Close()

	logger.Printf("Proxy started on %s forwarding to %s", listenAddr, targetAddr)

	connChan := make(chan net.Conn)

	for i := 0; i < runtime.NumCPU(); i++ {
		go handleConnections(connChan, targetAddr, logger)
	}

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			logger.Printf("Error accepting connection on %s: %v", listenAddr, err)
			continue
		}
		connChan <- clientConn
	}
}

// handleConnections sets up bidirectional copying between the client and remote server.
func handleConnections(connChan <-chan net.Conn, targetAddr string, logger *log.Logger) {
	for clientConn := range connChan {
		go func(conn net.Conn) {
			defer conn.Close()

			clientAddr := conn.RemoteAddr().String()
			logger.Printf("New connection: %s -> %s", clientAddr, targetAddr)

			serverConn, err := net.Dial("tcp", targetAddr)
			if err != nil {
				logger.Printf("Failed to connect to server %s: %v", targetAddr, err)
				return
			}
			defer serverConn.Close()

			done := make(chan struct{}, 2)

			go func() {
				_, err := io.Copy(serverConn, conn)
				if err != nil && err != io.EOF {
					logger.Printf("Error copying from client %s to server %s: %v", clientAddr, targetAddr, err)
				}
				done <- struct{}{}
			}()

			go func() {
				_, err := io.Copy(conn, serverConn)
				if err != nil && err != io.EOF {
					logger.Printf("Error copying from server %s to client %s: %v", targetAddr, clientAddr, err)
				}
				done <- struct{}{}
			}()

			<-done
			<-done

			logger.Printf("Connection closed: %s -> %s", clientAddr, targetAddr)
		}(clientConn)
	}
}
