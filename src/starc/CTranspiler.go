package main

import (
	"fmt"
    "os"
)

type Transpiler struct {
    fileName string
    parsedFile string
    current int
}

func (t *Transpiler) NewCFile() *os.File  {
    newFile, _ := os.Create(fmt.Sprintf("%s.c", t.fileName))
    stdioHeader := []byte("#include <stdio.h>")
    os.WriteFile(fmt.Sprintf("%s.c", t.fileName), stdioHeader, 0644)
    return newFile
}