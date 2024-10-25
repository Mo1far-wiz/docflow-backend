package docgenerator

import (
	"docflow-backend/models"
	"fmt"
	"os"

	"github.com/signintech/gopdf"
)

func GeneratePDF(doc models.Doc, user models.User) (*gopdf.GoPdf, error) {
	const fontSize = 25

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4Landscape}) // Page size: A4 Landscape
	pdf.AddPage()

	// Add the TTF font
	err := pdf.AddTTFFont("font", "./assets/fonts/KyivTypeSans-Medium2.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}

	// Set the font for the document with a larger size
	err = pdf.SetFont("font", "", 24) // Larger font size for better readability
	if err != nil {
		return nil, fmt.Errorf("failed to set font: %v", err)
	}

	// Center the document name and number
	pdf.SetX(0)  // Reset X to center
	pdf.SetY(40) // Set Y position for title
	pdf.CellWithOption(&gopdf.Rect{W: 842, H: 40}, fmt.Sprintf("%s\n (document ID №%d)", doc.DocName, doc.ID), gopdf.CellOption{Align: gopdf.Center})

	// Add a new line for spacing
	pdf.Br(30) // Increased space after title

	// Set font size for user information
	err = pdf.SetFont("font", "", fontSize) // Increase font size for user information
	if err != nil {
		return nil, fmt.Errorf("failed to set font: %v", err)
	}

	/// Create the content to center
	content := fmt.Sprintf(
		"Issued for %s %s, student of %s, %s, %d year of study. "+
			"This document serves as confirmation of the student's status at the institution.",
		user.FirstName, user.LastName, doc.Faculty, doc.Specialty, doc.YearOfStudy,
	)

	// Split content into lines
	lines := splitLines(content, 70) // Adjust max line length for A4 Landscape size

	// Calculate total height for centering
	lineHeight := fontSize*float64(len(lines)) + 25*float64(len(lines)-1) // Adjusted spacing
	totalPageHeight := 595.0                                              // A4 landscape height
	centerY := (totalPageHeight - lineHeight) / 2                         // Calculate Y position for vertical centering

	// Set X position for text to start from the center
	pdf.SetX(10)      // Set X position to left margin
	pdf.SetY(centerY) // Set calculated Y position

	// Add content with the user's information
	for _, line := range lines {
		err = pdf.Text(line)
		if err != nil {
			return nil, fmt.Errorf("failed to write text: %v", err)
		}
		pdf.Br(25) // Increased spacing between lines
	}

	// Add the creation date at the bottom left
	dateStr := doc.DateTime.Format("02.01.2006") // Format as dd.mm.yyyy
	pdf.SetY(540)                                // Adjust Y position for bottom left (just above bottom margin)
	pdf.SetX(10)                                 // Set X position to left
	pdf.SetFont("font", "", fontSize-6)          // Font size for the date
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
