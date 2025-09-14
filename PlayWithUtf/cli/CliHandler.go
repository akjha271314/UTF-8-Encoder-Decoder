package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"playwithutf/services"
	"strconv"
	"strings"
)

// StartInteractiveSession handles the main application loop, reading commands from stdin.
func StartInteractiveSession() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("UTF-8 Encoder/Decoder Interactive Mode")
	fmt.Println("Enter commands (e.g., --encode=\"Hello\" or --decode=\"48,65,6C\"). Press Ctrl+C to exit.")

	// The main application loop.
	for {
		fmt.Print(">> ") // Print a prompt for the user.

		// Wait for the user to enter a line.
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		// If the line is empty, skip to the next loop iteration.
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split the line into arguments, similar to how os.Args works.
		args := strings.Fields(line)

		// Process the arguments from the user's line.
		err := processArgs(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func processArgs(args []string) error {
	// Create a new FlagSet to handle the command-line flags.
	fs := flag.NewFlagSet("UTF-8", flag.ExitOnError)

	// Define flags for encode and decode operations.
	var encodeString string
	fs.StringVar(&encodeString, "encode", "", "The string to be UTF-8 encoded")

	var decodeHexBytes string
	fs.StringVar(&decodeHexBytes, "decode", "", "A comma-separated string of hexadecimal byte values to be UTF-8 decoded (e.g., F0,9F,98,8A)")

	// Parse the provided arguments.
	if err := fs.Parse(args); err != nil {
		return err
	}

	if encodeString != "" {
		// Encoding logic
		var allEncodedBytes []byte
		for _, r := range encodeString {
			encoded := services.Utf8Encode(r)
			allEncodedBytes = append(allEncodedBytes, encoded...)
		}
		fmt.Printf("Original string: %s\n", encodeString)
		fmt.Printf("Encoded byte slice: %v\n", allEncodedBytes)
		return nil
	}

	if decodeHexBytes != "" {
		// Decoding logic
		parts := strings.Split(decodeHexBytes, ",")
		var inputBytes []byte

		for _, part := range parts {
			// Trim spaces and convert from hex string to a byte.
			hexValue, err := strconv.ParseUint(strings.TrimSpace(part), 16, 8)
			if err != nil {
				return fmt.Errorf("invalid hexadecimal value provided: %s", part)
			}
			inputBytes = append(inputBytes, byte(hexValue))
		}

		fmt.Printf("Original byte slice (parsed from hex): %v\n", inputBytes)
		fmt.Printf("Decoded runes:\n")

		// Iterate over the byte slice, decoding one rune at a time.
		for i := 0; i < len(inputBytes); {
			r, size := services.Utf8Decode(inputBytes[i:])
			if size == 0 {
				fmt.Printf("Error: Invalid UTF-8 sequence at position %d\n", i)
				break
			}
			fmt.Printf(" - Rune: U+%X, Character: %c, Size: %d bytes\n", r, r, size)
			i += size // Advance the counter by the number of bytes consumed.
		}
		return nil
	}

	// If no flags were provided, show usage.
	return fmt.Errorf("please provide either --encode or --decode flag")
}
