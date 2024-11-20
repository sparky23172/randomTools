package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const bufferSize = 10 * 1024 * 1024 // 10 MB

var logger *zap.Logger

func main() {
	// Define flags
	fileFlag := flag.String("file", "", "File to scan (optional)")
	dirFlag := flag.String("dir", "", "Directory to scan (optional)")
	regexFlag := flag.String("regex", "", "Custom regex pattern to match (optional)")
	minCharsFlag := flag.Int("minchars", 4, "Minimum number of characters in encoded string to check")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Initialize logger based on debug flag
	var err error
	if *debugFlag {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Add colors for better visibility
		logger, err = cfg.Build()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		fmt.Printf("Failed to initialize logger: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Validate flags
	if *fileFlag == "" && *dirFlag == "" {
		logger.Error("Either -file or -dir must be specified.")
		os.Exit(1)
	}

	// Create the regex dynamically based on the minimum characters
	var regexPattern string
	if *regexFlag != "" {
		regexPattern = *regexFlag
	} else {
		// Base64 regex with customizable minimum length
		minChars := *minCharsFlag
		regexPattern = fmt.Sprintf(`(([A-Za-z0-9+\/]{%d,})([A-Za-z0-9+\/]{4}|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{2}==))`, minChars)
	}

	// Compile the regex pattern
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		logger.Error("Invalid regex pattern", zap.Error(err))
		os.Exit(1)
	}

	if *fileFlag != "" {
		processFile(*fileFlag, re, *debugFlag)
	} else if *dirFlag != "" {
		err := filepath.Walk(*dirFlag, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logger.Warn("Error accessing file", zap.String("path", path), zap.Error(err))
				return nil
			}
			if !info.IsDir() {
				processFile(path, re, *debugFlag)
			}
			return nil
		})
		if err != nil {
			logger.Error("Error walking directory", zap.Error(err))
		}
	}
}

func processFile(fileName string, re *regexp.Regexp, debugMode bool) {
	if debugMode {
		logger.Debug("Scanning file", zap.String("file", fileName))
	}

	file, err := os.Open(fileName)
	if err != nil {
		logger.Warn("Failed to open file", zap.String("file", fileName), zap.Error(err))
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, bufferSize)
	lineNumber := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				logger.Warn("Error reading file", zap.String("file", fileName), zap.Error(err))
			}
			break
		}

		lineNumber++
		// Find matches using the regex
		matches := re.FindAllString(line, -1)
		for _, match := range matches {
			// Attempt to decode the match as base64
			decoded, err := base64.StdEncoding.DecodeString(match)
			if err != nil {
				logger.Debug("Failed to decode base64", zap.String("encoded", match), zap.Error(err))
				continue
			}

			// Check if decoded string is valid UTF-8
			if !utf8.Valid(decoded) {
				logger.Debug("Invalid UTF-8 decoded string", zap.String("decoded", string(decoded)))
				continue
			}

			// Ensure the decoded string is printable
			if !isPrintableUTF8(decoded) {
				logger.Debug("Decoded string contains non-printable characters", zap.String("decoded", string(decoded)))
				continue
			}

			logger.Info("Valid base64 string found",
				zap.String("file", fileName),
				zap.Int("line", lineNumber),
				zap.String("encoded", match),
				zap.String("decoded", string(decoded)))
		}
	}
}

func isPrintableUTF8(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		if r == utf8.RuneError || !isPrintableRune(r) {
			return false
		}
		data = data[size:]
	}
	return true
}

func isPrintableRune(r rune) bool {
	// Check if the rune is a printable character
	return (r >= ' ' && r <= '~') || r == '\n' || r == '\t'
}
