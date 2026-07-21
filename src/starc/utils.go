package main;

import (
    "os"
    "strings"
    "path/filepath"
)

type Color string

const (
    Red = "\033[31m"
    Orange = "\033208m"
    Green = "\033[32m"
    Yellow = "\033[33m"
    BrightYellow = "\03393m"
    Blue = "\033[34m"
    Purple = "\033[35m"
    BrightPurple = "\033[95m"
    Reset = "\033[0m"
)

func DirFiles(dir string) []string {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil
    }
    var files []string
    for _, e := range entries {
        if !e.IsDir() {
            files = append(files, e.Name())
        }
    }
    return files
}

func SearchProjectRoot(marker string) string {
    dir, err := os.Getwd()
    if err != nil {
        return "."
    }
    
    for {
        markerPath := filepath.Join(dir, marker)
        if _, err := os.Stat(markerPath); err == nil {
            return dir
        }
        
        parent := filepath.Dir(dir)
        if parent == dir {
            return "."
        }
        dir = parent
    }
}

func TrimBefore(str string, trim string) string {
    var res string = str
    last := strings.LastIndex(str, trim)
    if last != -1 {res = str[last+1:]}
    return res
}

func MatchStr(str string, compare ...string) bool {
    for _, comp := range compare {
        if str == comp {return true}
    }
    return false
}