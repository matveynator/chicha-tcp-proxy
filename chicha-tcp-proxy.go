// Package main provides a TCP proxy that listens on specified local ports
// and forwards incoming traffic to corresponding remote destinations.
// This implementation includes log rotation, file compression of old logs,
// and uses channels and goroutines to handle concurrent connections.
// No mutexes are used; concurrency is managed by channels and goroutines,
// adhering to Go proverbs and common Go idioms.
package main

import (
	"compress/gzip" // Package gzip provides support for reading and writing GZIP format compressed files.
	"flag"          // Package flag implements command-line flag parsing.
	"fmt"           // Package fmt implements formatted I/O.
	"io"            // Package io provides basic interfaces to I/O primitives.
	"log"           // Package log implements a simple logging package.
	"net"           // Package net provides a portable interface for network I/O.
	"os"            // Package os provides a platform-independent interface to operating system functionality.
	"runtime"       // Package runtime provides operations that interact with Go's runtime system.
	"strings"       // Package strings implements simple functions to manipulate UTF-8 encoded strings.
	"time"          // Package time provides functionality for measuring and displaying time.
)

// Route describes a single forwarding route configuration from a local port to a remote address.
// LocalPort: The local port on which the proxy listens (e.g. "8080")
// RemoteIP: The target server IP address to forward traffic to (e.g. "46.4.70.114")
// RemotePort: The remote port on the target server to forward traffic to (e.g. "80")
type Route struct {
	LocalPort  string // The local port number as a string.
	RemoteIP   string // The remote IP address as a string.
	RemotePort string // The remote port number as a string.
}

func main() {
	// routesFlag holds the comma-separated list of routes in the format LOCALPORT:REMOTEIP:REMOTEPORT
	routesFlag := flag.String("routes", "", "Comma-separated list of routes in the format LOCALPORT:REMOTEIP:REMOTEPORT")
	// logFile specifies the path to the log file where proxy activity will be logged.
	logFile := flag.String("log", "chicha-tcp-proxy.log", "Path to the log file")
	// rotationFrequency specifies how often the log file should be rotated.
	rotationFrequency := flag.Duration("rotation", 24*time.Hour, "Log rotation frequency (e.g. 24h, 1h, etc.)")

	// Parse the provided command-line flags.
	flag.Parse()

	// Validate that the required routes flag is provided.
	if *routesFlag == "" {
		log.Fatal("Error: The -routes flag is required.")
	}

	// Parse the routes from the provided string.
	routes, err := parseRoutes(*routesFlag)
	if err != nil {
		log.Fatalf("Error parsing routes: %v", err)
	}
	if len(routes) == 0 {
		log.Fatalf("Error: no valid routes found in '%s'", *routesFlag)
	}

	// Print basic startup information: routes, log file, and rotation frequency.
	fmt.Println("========== CHICHA TCP PROXY ==========")
	fmt.Println("Routes:")
	for _, route := range routes {
		fmt.Printf("  LocalPort=%s -> RemoteIP=%s RemotePort=%s\n", route.LocalPort, route.RemoteIP, route.RemotePort)
	}
	fmt.Printf("Log file: %s\n", *logFile)
	fmt.Printf("Log rotation frequency: %v\n", *rotationFrequency)
	fmt.Println("======================================")

	// Set up the logger that will write to the specified log file.
	logger, file, err := setupLogger(*logFile)
	if err != nil {
		log.Fatalf("Error setting up logger: %v", err)
	}

	log.Printf("Starting chicha-tcp-proxy")

	// Set the number of OS threads to use based on the number of CPUs available.
	// According to Go proverbs, "Don't communicate by sharing memory; share memory by communicating."
	// By default Go does this well, but we explicitly set it for clarity.
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
	logger.Printf("Using %d CPU cores", numCPUs)
	log.Printf("Using %d CPU cores", numCPUs)

	// Start the log rotation in a separate goroutine. This periodically rotates logs without blocking the main execution.
	go rotateLogs(*logFile, file, logger, *rotationFrequency)

	// Start the proxy servers for each route in separate goroutines. This allows concurrent handling of multiple routes.
	for _, route := range routes {
		// Inform about starting a proxy instance for this route.
		logger.Printf("Starting proxy for route: local=%s remote=%s:%s", route.LocalPort, route.RemoteIP, route.RemotePort)
		log.Printf("Starting proxy for route: local=%s remote=%s:%s", route.LocalPort, route.RemoteIP, route.RemotePort)

		// Launch a goroutine to handle incoming connections on the specified local port and forward them to the remote address.
		go startProxy(":"+route.LocalPort, route.RemoteIP+":"+route.RemotePort, logger)
	}

	// Block indefinitely to keep the main function running.
	// Using select{} is a common idiom for blocking forever.
	select {}
}

// parseRoutes parses a comma-separated string of route definitions into a slice of Route objects.
// The expected format for each route is "LOCALPORT:REMOTEIP:REMOTEPORT".
func parseRoutes(routesFlag string) ([]Route, error) {
	// Split the input on commas to separate individual route definitions.
	parts := strings.Split(routesFlag, ",")
	var routes []Route
	for _, part := range parts {
		// Split each route into its components: local port, remote IP, and remote port.
		segments := strings.Split(part, ":")
		if len(segments) != 3 {
			return nil, fmt.Errorf("invalid route format: '%s' (expected LOCALPORT:REMOTEIP:REMOTEPORT)", part)
		}
		// Construct a Route struct and add it to the routes slice.
		routes = append(routes, Route{
			LocalPort:  segments[0],
			RemoteIP:   segments[1],
			RemotePort: segments[2],
		})
	}
	return routes, nil
}

// setupLogger creates or opens the specified log file and returns a logger and the file handle.
// If the file does not exist, it will be created. Logs are appended if the file already exists.
func setupLogger(logFile string) (*log.Logger, *os.File, error) {
	// Open or create the log file with append mode.
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file '%s': %v", logFile, err)
	}

	// Create a new logger that writes to the opened file.
	logger := log.New(file, "", log.LstdFlags)
	return logger, file, nil
}

// rotateLogs handles periodic rotation of the current log file. After rotation,
// it compresses the old log file and starts a new one.
// This runs indefinitely in a goroutine.
func rotateLogs(logFile string, file *os.File, logger *log.Logger, frequency time.Duration) {
	for {
		// Sleep for the specified rotation frequency before rotating logs again.
		time.Sleep(frequency)

		// Close the current log file before renaming.
		file.Close()

		// Create a rotated filename based on the current date.
		rotatedFile := logFile + "." + time.Now().Format("2006-01-02")
		if err := os.Rename(logFile, rotatedFile); err != nil {
			// If renaming fails, log the error and attempt to reopen the current log file to continue logging.
			logger.Printf("Error rotating logs: %v", err)
			log.Printf("Error rotating logs: %v", err)

			newFile, err2 := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err2 != nil {
				// If reopening also fails, we must terminate since we have nowhere to log.
				logger.Fatalf("Failed to reopen log file after rotation error: %v", err2)
				log.Fatalf("Failed to reopen log file after rotation error: %v", err2)
			}
			file = newFile
			logger.SetOutput(file)
			continue
		}

		// After successful rename, open a new log file with the original name to continue logging.
		newFile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatalf("Failed to create new log file after rotation: %v", err)
			log.Fatalf("Failed to create new log file after rotation: %v", err)
		}
		file = newFile
		logger.SetOutput(file)

		// Inform that the log was rotated and now compressing the old log file.
		logger.Println("Log file rotated successfully, now compressing old log...")
		log.Println("Log file rotated successfully, now compressing old log...")

		// Compress the old log file to save space and then remove the uncompressed version.
		if err := compressFile(rotatedFile); err != nil {
			logger.Printf("Error compressing rotated file: %v", err)
			log.Printf("Error compressing rotated file: %v", err)
		} else {
			logger.Printf("Compression successful: %s.gz", rotatedFile)
			log.Printf("Compression successful: %s.gz", rotatedFile)
			if err := os.Remove(rotatedFile); err != nil {
				logger.Printf("Error removing uncompressed rotated file: %v", err)
				log.Printf("Error removing uncompressed rotated file: %v", err)
			}
		}
	}
}

// compressFile takes a filename and compresses it using gzip, creating a .gz file.
// After compression, the original file can be removed by the caller.
func compressFile(filename string) error {
	// Open the original file for reading.
	original, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file for compression: %v", err)
	}
	defer original.Close()

	// Create a new .gz file for writing the compressed data.
	gzFile, err := os.OpenFile(filename+".gz", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create gz file: %v", err)
	}
	defer gzFile.Close()

	// Create a gzip writer to compress data as it's written.
	gzWriter := gzip.NewWriter(gzFile)
	defer gzWriter.Close()

	// Copy all data from the original file into the gzip writer (which compresses it).
	if _, err := io.Copy(gzWriter, original); err != nil {
		return fmt.Errorf("failed to copy data for compression: %v", err)
	}

	return nil
}

// startProxy listens on the specified local address and forwards all connections to the target address.
// Each accepted connection is passed through a channel to be handled by worker goroutines.
func startProxy(listenAddr, targetAddr string, logger *log.Logger) {
	// Listen for incoming TCP connections on the given local address.
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("Failed to start proxy on %s: %v", listenAddr, err)
	}
	defer listener.Close()

	logger.Printf("Proxy started on %s forwarding to %s", listenAddr, targetAddr)

	// Create a channel to distribute accepted connections to worker goroutines.
	connChan := make(chan net.Conn)

	// Spawn worker goroutines to handle connections concurrently.
	// Using the number of CPUs as the number of workers is a common approach.
	for i := 0; i < runtime.NumCPU(); i++ {
		go handleConnections(connChan, targetAddr, logger)
	}

	// Continuously accept new client connections and send them to the channel for processing.
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			logger.Printf("Error accepting connection on %s: %v", listenAddr, err)
			continue // If there's a transient error, continue accepting the next connection.
		}
		// Send the new connection to one of the worker goroutines via the channel.
		connChan <- clientConn
	}
}

// handleConnections receives connections from the channel and sets up bidirectional data copying
// between the client and the remote server. It uses goroutines for each direction of traffic.
// This function blocks, reading from connChan, until the channel is closed or the program is terminated.
func handleConnections(connChan <-chan net.Conn, targetAddr string, logger *log.Logger) {
	for {
		// Use a select statement for possible future expansions (like graceful shutdown).
		select {
		case clientConn, ok := <-connChan:
			if !ok {
				// If the channel is closed, return to stop this worker.
				return
			}

			// For each client connection, start a new goroutine to handle forwarding.
			go func(conn net.Conn) {
				defer conn.Close()

				clientAddr := conn.RemoteAddr().String()
				logger.Printf("New connection: %s -> %s", clientAddr, targetAddr)

				// Dial the target server.
				serverConn, err := net.Dial("tcp", targetAddr)
				if err != nil {
					logger.Printf("Failed to connect to server %s: %v", targetAddr, err)
					return
				}
				defer serverConn.Close()

				// done channel signals when copying in each direction finishes.
				done := make(chan struct{}, 2)

				// Copy data from client to server.
				go func() {
					_, err := io.Copy(serverConn, conn)
					if err != nil && err != io.EOF {
						logger.Printf("Error copying from client %s to server %s: %v", clientAddr, targetAddr, err)
					}
					done <- struct{}{}
				}()

				// Copy data from server to client.
				go func() {
					_, err := io.Copy(conn, serverConn)
					if err != nil && err != io.EOF {
						logger.Printf("Error copying from server %s to client %s: %v", targetAddr, clientAddr, err)
					}
					done <- struct{}{}
				}()

				// Wait for both copy operations to complete before closing the connection.
				<-done
				<-done

				logger.Printf("Connection closed: %s -> %s", clientAddr, targetAddr)
			}(clientConn)
		}
	}
}
