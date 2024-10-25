package docgenerator

import (
	"docflow-backend/models"
	"fmt"
	"os"

	"github.com/signintech/gopdf"
)

var docContents = map[string]string{
	"Certificate of study":                                  "This document serves as confirmation of the student's status at the institution.",
	"Certificate of Tuition Fees":                           "This is to certify that the student has fulfilled the payment obligations for tuition fees as required by the institution. The tuition fees have been paid in full for the relevant academic period, ensuring the student’s enrollment and participation in the designated program.",
	"Certificate of fulfillment of the Corporate Agreement": "This is to certify that the obligations set forth under the Corporate Agreement between University and Student have been successfully fulfilled in accordance with the agreed terms. Both parties have met their responsibilities, completing all activities and commitments specified within the agreement to mutual satisfaction.",
	"Certificate of storage of original documents":          "This is to certify that the original documents submitted by the student are securely stored by University. These documents have been received, verified, and retained in accordance with the university's policies and procedures.",
	"Certificate of payment for the contract":               "This is to certify that the payment required under the contract between University and the student has been received in full. The financial obligations outlined in the agreement have been fulfilled in accordance with the terms and conditions specified.",
}

func GeneratePDF(doc models.Doc, user models.User) (*gopdf.GoPdf, error) {
	contents, exists := docContents[doc.DocName]
	if !exists {
		return nil, fmt.Errorf("no such certificate exists")
	}

	const fontSize = 22

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4Landscape})
	pdf.AddPage()

	err := pdf.AddTTFFont("font", "./assets/fonts/KyivTypeSans-Medium2.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}

	err = pdf.SetFont("font", "", fontSize)
	if err != nil {
		return nil, fmt.Errorf("failed to set font: %v", err)
	}

	pdf.SetX(0)
	pdf.SetY(40)
	pdf.CellWithOption(&gopdf.Rect{W: 842, H: 40}, fmt.Sprintf("%s\n (document ID №%d)", doc.DocName, doc.ID), gopdf.CellOption{Align: gopdf.Center})

	pdf.Br(30)

	err = pdf.SetFont("font", "", fontSize)
	if err != nil {
		return nil, fmt.Errorf("failed to set font: %v", err)
	}

	content := fmt.Sprintf(
		"Issued for %s %s, student of %s, %s, %d year of study. "+
			contents,
		user.FirstName, user.LastName, doc.Faculty, doc.Specialty, doc.YearOfStudy,
	)

	lines := splitLines(content, 65)

	lineHeight := fontSize*float64(len(lines)) + 25*float64(len(lines)-1)
	totalPageHeight := 595.0
	centerY := (totalPageHeight - lineHeight) / 2

	pdf.SetX(15)
	pdf.SetY(centerY)

	for _, line := range lines {
		err = pdf.Text(line)
		if err != nil {
			return nil, fmt.Errorf("failed to write text: %v", err)
		}
		pdf.Br(25)
	}

	dateStr := doc.DateTime.Format("02.01.2006")
	pdf.SetY(540)
	pdf.SetX(15)
	pdf.SetFont("font", "", fontSize-6)
	pdf.Cell(nil, "Issued at : "+dateStr)

	err = addLogo(pdf, "./assets/logo_stamp.png")
	if err != nil {
		return nil, fmt.Errorf("failed to add logo: %v", err)
	}

	return pdf, nil
}

// Helper function to split long text into smaller lines.
func splitLines(text string, maxLen int) []string {
	var lines []string
	for len(text) > maxLen {
		splitAt := maxLen
		for i := maxLen; i >= 0; i-- {
			if text[i] == ' ' {
				splitAt = i
				break
			}
		}
		lines = append(lines, text[:splitAt])
		text = text[splitAt+1:] // Skip the space
	}
	lines = append(lines, text) // Add the remaining part
	return lines
}

func addLogo(pdf *gopdf.GoPdf, imagePath string) error {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("logo not found: %v", err)
	}

	// Calculate position: bottom-right corner
	pageWidth, pageHeight := 842.0, 595.0   // A4 Landscape dimensions
	imageWidth, imageHeight := 100.0, 100.0 // Example dimensions for the logo

	x := pageWidth - imageWidth - 50
	y := pageHeight - imageHeight - 50

	// Add the image to the PDF
	err := pdf.Image(imagePath, x, y, &gopdf.Rect{W: imageWidth, H: imageHeight})
	if err != nil {
		return fmt.Errorf("failed to add image: %v", err)
	}
	return nil
}
