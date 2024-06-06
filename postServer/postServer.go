package main

import (
    "encoding/base64"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
)

func postHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Unable to read request body", http.StatusBadRequest)
            return
        }
        // Log the POST data to the server console
        log.Printf("Received POST body: %s", string(body))
        fmt.Fprintf(w, "Received POST body: %s", string(body))
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func b64Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Unable to read request body", http.StatusBadRequest)
            return
        }
        decoded, err := base64.StdEncoding.DecodeString(string(body))
        if err != nil {
            errorMessage := fmt.Sprintf("Base64 Decoding Failed: '%s'", string(body))
            log.Printf("Base64 Decoding Failed: '%s'", string(body))
            http.Error(w, errorMessage, http.StatusBadRequest)
            return
        }
        log.Printf("Decoded data: %s", string(decoded))
        fmt.Fprintf(w, "Decoded data: %s", string(decoded))
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func printNetworkInterfaces() {
    interfaces, err := net.Interfaces()
    if err != nil {
        log.Fatalf("Failed to get network interfaces: %v", err)
    }

    for _, iface := range interfaces {
        addrs, err := iface.Addrs()
        if err != nil {
            log.Printf("Failed to get addresses for interface %v: %v", iface.Name, err)
            continue
        }

        for _, addr := range addrs {
            fmt.Printf("Interface: %v, IP Address: %v\n", iface.Name, addr.String())
        }
    }
}

func serveHTTPServer(port string) {
    http.HandleFunc("/", postHandler)
    http.HandleFunc("/b64", b64Handler)

    printNetworkInterfaces()

    fmt.Printf("Starting server at :%s\n", port)
    log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func base64EncodeFile(filePath string) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        log.Fatalf("Failed to read file: %v", err)
    }
    encoded := base64.StdEncoding.EncodeToString(data)
    fmt.Println(encoded)
}

func base64DecodeFile(filePath string) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        log.Fatalf("Failed to read file: %v", err)
    }
    decoded, err := base64.StdEncoding.DecodeString(string(data))
    if err != nil {
        log.Fatalf("Failed to decode base64 data: %v", err)
    }
    fmt.Println(string(decoded))
}

func main() {
    encodeFilePath := flag.String("encode", "", "File to base64 encode")
    decodeFilePath := flag.String("decode", "", "File to base64 decode")
    port := flag.String("port", "8080", "Port to run the HTTP server on")
    flag.Parse()

    if *encodeFilePath != "" {
        base64EncodeFile(*encodeFilePath)
    } else if *decodeFilePath != "" {
        base64DecodeFile(*decodeFilePath)
    } else {
        serveHTTPServer(*port)
    }
}
