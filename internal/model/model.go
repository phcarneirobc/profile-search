package model

import "time"

type ProfileInfo struct {
    Name      string
    Location  string
    Bio       string
    Followers string
    Following string
    JoinDate  string
    Website   string
}

type Result struct {
    Platform     string
    URL          string
    Exists       bool
    Error        string
    ResponseTime time.Duration
    Info         *ProfileInfo
}