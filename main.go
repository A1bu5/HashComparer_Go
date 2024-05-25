package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// computeHash computes the MD5 and SHA256 hashes of a file
func computeHash(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	md5Hash := md5.New()
	sha256Hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(md5Hash, sha256Hash), file); err != nil {
		return "", "", err
	}

	md5Sum := hex.EncodeToString(md5Hash.Sum(nil))
	sha256Sum := hex.EncodeToString(sha256Hash.Sum(nil))

	return md5Sum, sha256Sum, nil
}

func main() {
	a := app.New()
	w := a.NewWindow("Hash Calculator")

	file1Entry := widget.NewEntry()
	file1Entry.SetPlaceHolder("Select file 1")

	file2Entry := widget.NewEntry()
	file2Entry.SetPlaceHolder("Select file 2")

	md5Label1 := widget.NewLabel("MD5 (File 1):")
	sha256Label1 := widget.NewLabel("SHA256 (File 1):")
	md5Label2 := widget.NewLabel("MD5 (File 2):")
	sha256Label2 := widget.NewLabel("SHA256 (File 2):")

	statusLabel := canvas.NewText("", theme.TextColor())
	statusLabel.TextSize = 16
	statusLabel.Alignment = fyne.TextAlignCenter
	statusLabel.Hide()

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	selectFile1 := widget.NewButton("Select File 1", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				file1Entry.SetText(reader.URI().Path())
			}
		}, w)
	})

	selectFile2 := widget.NewButton("Select File 2", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				file2Entry.SetText(reader.URI().Path())
			}
		}, w)
	})

	computeButton := widget.NewButton("Compute Hash", func() {
		filename1 := file1Entry.Text
		filename2 := file2Entry.Text

		if filename1 == "" && filename2 == "" {
			dialog.ShowError(fmt.Errorf("no file selected"), w)
			return
		}

		if filename1 != "" {
			md5Sum, sha256Sum, err := computeHash(filename1)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			md5Label1.SetText(fmt.Sprintf("MD5 (File 1): %s", md5Sum))
			sha256Label1.SetText(fmt.Sprintf("SHA256 (File 1): %s", sha256Sum))
		}

		if filename2 != "" {
			md5Sum, sha256Sum, err := computeHash(filename2)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			md5Label2.SetText(fmt.Sprintf("MD5 (File 2): %s", md5Sum))
			sha256Label2.SetText(fmt.Sprintf("SHA256 (File 2): %s", sha256Sum))
		}
	})

	compareButton := widget.NewButton("Compare Files", func() {
		file1 := file1Entry.Text
		file2 := file2Entry.Text

		if file1 == "" && file2 == "" {
			dialog.ShowError(fmt.Errorf("no files selected"), w)
			return
		}

		progressBar.Show()

		go func() {
			defer progressBar.Hide()

			time.Sleep(1 * time.Second) // Simulate some work being done

			var md5_1, sha256_1, md5_2, sha256_2 string
			var err1, err2 error

			if file1 != "" {
				md5_1, sha256_1, err1 = computeHash(file1)
			}
			if file2 != "" {
				md5_2, sha256_2, err2 = computeHash(file2)
			}

			if file1 != "" && err1 == nil {
				md5Label1.SetText(fmt.Sprintf("MD5 (File 1): %s", md5_1))
				sha256Label1.SetText(fmt.Sprintf("SHA256 (File 1): %s", sha256_1))
			}
			if file2 != "" && err2 == nil {
				md5Label2.SetText(fmt.Sprintf("MD5 (File 2): %s", md5_2))
				sha256Label2.SetText(fmt.Sprintf("SHA256 (File 2): %s", sha256_2))
			}

			statusLabel.Show()
			if file1 != "" && file2 != "" {
				areEqual := md5_1 == md5_2 && sha256_1 == sha256_2
				if areEqual {
					statusLabel.Text = "The Files Are The Same"
					statusLabel.Color = theme.PrimaryColor()
				} else {
					statusLabel.Text = "The Files Are Different"
					statusLabel.Color = theme.ErrorColor()
				}
			} else {
				statusLabel.Text = "Comparison requires two files"
				statusLabel.Color = theme.WarningColor()
			}

			statusLabel.Refresh()
		}()
	})

	content := container.NewVBox(
		file1Entry,
		selectFile1,
		file2Entry,
		selectFile2,
		computeButton,
		md5Label1,
		sha256Label1,
		md5Label2,
		sha256Label2,
		compareButton,
		progressBar,
		statusLabel,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 500))
	w.SetFixedSize(true)
	w.ShowAndRun()
}
