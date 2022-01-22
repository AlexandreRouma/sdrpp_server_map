package main

import (
    "io"
    "encoding/json"
    "context"
    "os"
    "fmt"
    "net/http"
    "github.com/jackc/pgx/v4"
)

var db *pgx.Conn

type DeviceInfo struct {
    Ip          string  `json:"ip"`        
    Port        int     `json:"port"`
    Sources     string  `json:"sources"`
    Antenna     string  `json:"antenna"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    Altitude    float64 `json:"altitude"`
}

type DevicesResponse struct {
    Error   string          `json:"error"`        
    Devices []DeviceInfo    `json:"devices"`
}

func apiDevicesHandler(w http.ResponseWriter, req *http.Request) {
    // Default response
    resp := DevicesResponse {
        Error: "",
    };

    // Query data
    rows, err := db.Query(context.Background(), "SELECT * FROM devices")
    if (err != nil) {
        resp.Error = "DATABASE_ERROR"

        // Send reply
        jsonResp, _ := json.Marshal(resp)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, string(jsonResp))
        return
    }

    // Read all rows
    for rows.Next() {
        dev := DeviceInfo {}
        err = rows.Scan(&dev.Ip, &dev.Port, &dev.Sources, &dev.Antenna, &dev.Latitude, &dev.Longitude, &dev.Altitude)
        if (err != nil) {
            resp.Error = "DATABASE_ERROR"

            // Send reply
            jsonResp, _ := json.Marshal(resp)
            w.Header().Set("Content-Type", "application/json")
            fmt.Fprint(w, string(jsonResp))
            return
        }
        resp.Devices = append(resp.Devices, dev)
    }
    rows.Close()

    // Send reply
    jsonResp, _ := json.Marshal(resp)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(jsonResp))
}

type RegisterResponse struct {
    Error          string       `json:"error"`
}

func apiRegister(w http.ResponseWriter, req *http.Request) {
    // Default response
    resp := RegisterResponse {
        Error: "",
    }
    
    // Decode the body
    _body, err := io.ReadAll(req.Body);
    if (err != nil) {
        resp.Error = "INVALID_REQUEST_BODY"

        // Send reply
        jsonResp, _ := json.Marshal(resp)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, string(jsonResp))
    }

    resp.Error = string(_body)

    // Send reply
    jsonResp, _ := json.Marshal(resp)
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(jsonResp))

    // Log addition
    fmt.Println("Registered new device")
}

func main() {
    // Connect to database
    var err error
    db, err = pgx.Connect(context.Background(), "postgres://postgres:12345678@localhost:5432/sdrpp_server_map")
    if (err != nil) {
        fmt.Fprintf(os.Stderr, "Could not connect to database: %v\n", err)
        os.Exit(-1)
    }
    defer db.Close(context.Background())

    // Static files
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)

    // API Handling
    http.HandleFunc("/api/devices", apiDevicesHandler)
    http.HandleFunc("/api/register", apiRegister)

    // Listen for connections
    fmt.Printf("Ready, listening on TODO:TODO\n")
    http.ListenAndServe(":8090", nil)
}