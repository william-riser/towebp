package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chai2010/webp"
)

var supportedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func convert(path string, quality float32) error {
	in, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer in.Close()

	img, _, err := image.Decode(in)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	outPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".webp"
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer out.Close()

	if err := webp.Encode(out, img, &webp.Options{Quality: quality}); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

func main() {
	dir := flag.String("dir", ".", "directory containing images (recursed)")
	quality := flag.Float64("quality", 80, "webp quality 0-100")
	workers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	flag.Parse()

	start := time.Now()

	var paths []string
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if supportedExts[strings.ToLower(filepath.Ext(path))] {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("walk: %v", err)
	}

	total := len(paths)
	if total == 0 {
		fmt.Println("No supported images found.")
		return
	}

	fmt.Printf("Converting %d images with %d workers (q=%.0f)...\n", total, *workers, *quality)

	jobs := make(chan string)
	var wg sync.WaitGroup
	var ok, fail int64

	for range *workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if err := convert(p, float32(*quality)); err != nil {
					atomic.AddInt64(&fail, 1)
					log.Printf("FAIL %s: %v", p, err)
				} else {
					atomic.AddInt64(&ok, 1)
				}
			}
		}()
	}

	for _, p := range paths {
		jobs <- p
	}
	close(jobs)
	wg.Wait()

	fmt.Printf("Done in %v: %d succeeded, %d failed\n",
		time.Since(start).Round(time.Millisecond), ok, fail)
}