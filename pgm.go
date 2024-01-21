package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// pgm représente une image pgm avec des données uint8, largeur, hauteur, magicNumber et une valeur maximale.
type pgm struct {
	data          [][]uint8 // Données de l'image
	width, height int       // Dimensions de l'image (largeur et hauteur)
	magicNumber   string    // Numéro magique du format pgm (P2 ou P4)
	max           uint8     // Valeur maximale de l'image
}

// Readpgm lit un fichier pgm et retourne un objet pgm rempli avec les données du fichier.
func ReadPGM(filename string) (*pgm, error) {
	var mypgmR pgm
	var mypgm *pgm = &mypgmR
	var caractere []rune
	var verifMN bool = false
	var verifDim bool = false
	var verifmax bool = false
	k := 0

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		caractere = []rune(line)
		currentLine := 0
		if string(caractere[0]) == "#" {
			continue // Ignorer les commentaires
		}
		if !verifMN {
			mypgm.magicNumber = line
			verifMN = true
			continue
		}

		if !verifDim {
			curLine := strings.Fields(line)
			mypgm.height, err = strconv.Atoi(curLine[0])
			if err != nil {
				fmt.Println("erreur dans la lecture de la longueur")
			}
			mypgm.width, err = strconv.Atoi(curLine[1])
			if err != nil {
				fmt.Println("erreur dans la lecture de la largeur")
			}
			mypgm.data = make([][]uint8, mypgm.height)
			for k := range mypgm.data {
				mypgm.data[k] = make([]uint8, mypgm.width)
			}
			verifDim = true
			continue
		}

		if !verifmax {
			resultUint64, err := strconv.ParseUint(line, 10, 64)
			result := uint8(resultUint64)
			mypgm.max = result

			if err != nil {
				fmt.Println("Error in the format:", err)
				break
			}
			verifmax = true
			continue
		}

		if mypgm.magicNumber == "P2" {
			words := strings.Fields(line)
			for i := range mypgm.data[currentLine] {
				curDat, err := strconv.Atoi(words[i])
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
					break
				}
				mypgm.data[k][i] = uint8(curDat)
			}
			k++
		}

		if mypgm.magicNumber == "P4" {
			for k := 0; k < mypgm.height; k++ {
				ligneInt := strings.Split(line, " ")
				for i := 0; i < len((ligneInt)); i++ {
					intStatutVerif, err := strconv.Atoi(ligneInt[i])
					if err != nil {
						fmt.Println("Error converting string to integer:", err)
						break
					}
					if intStatutVerif > int(mypgm.max) || intStatutVerif < 0 {
						fmt.Println("format invalide")
						break
					}
					mypgm.data[k][i] = uint8(intStatutVerif)
				}
			}
		}
	}
	return mypgm, nil
}

func (pgm *pgm) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écriture de l'en-tête
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Écriture des données de l'image
	for _, row := range pgm.data {
		for _, pixel := range row {
			fmt.Fprintf(file, "%d ", pixel)
		}
		fmt.Fprintln(file) // Nouvelle ligne après chaque ligne d'image
	}

	return nil
}

// Size retourne les dimensions de l'image (largeur et hauteur).
func (pgm *pgm) Size() (int, int) {
	return pgm.height, pgm.width
}

// At retourne la valeur à la position (x, y) dans les données de l'image.
func (pgm *pgm) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set modifie la valeur à la position (x, y) dans les données de l'image.
func (pgm *pgm) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Invert inverse toutes les valeurs de l'image.
func (pgm *pgm) Invert() {
	for i := range pgm.data {
		for k := range pgm.data[i] {
			pgm.data[i][k] = uint8(int(pgm.max) - int(pgm.data[i][k]))
		}
	}
}

// Flip effectue une opération de retournement horizontal sur l'image.
func (pgm *pgm) Flip() {
	if pgm.width <= 0 || pgm.height <= 0 {
		fmt.Println("Invalid dimensions")
		return
	}

	sliceUpdate := make([][]uint8, pgm.height)
	for i := range sliceUpdate {
		sliceUpdate[i] = make([]uint8, pgm.width)
	}

	for i := 0; i < pgm.height; i++ {
		for k := 0; k < pgm.width; k++ {
			sliceUpdate[i][k] = pgm.data[i][pgm.width-k-1]
		}
	}
	pgm.data = sliceUpdate
}

// Flop effectue une opération de retournement vertical sur l'image.
func (pgm *pgm) Flop() {
	if pgm.width <= 0 || pgm.height <= 0 {
		fmt.Println("Invalid dimensions")
		return
	}

	sliceUpdate := make([][]uint8, pgm.height)
	for i := range sliceUpdate {
		sliceUpdate[i] = make([]uint8, pgm.width)
	}

	for i := 1; i <= pgm.height; i++ {
		for k := 0; k < pgm.width; k++ {
			sliceUpdate[i-1][k] = pgm.data[pgm.height-i][k]
		}
	}
	pgm.data = sliceUpdate
}

// SetMagicNumber modifie le MagicNumber
func (pgm *pgm) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image pgm.
func (pgm *pgm) SetMaxValue(maxValue uint8) {
	// Sauvegarder la valeur maximale d'origine
	oldMax := pgm.max

	// Mettre à jour la nouvelle valeur maximale
	pgm.max = maxValue

	// Ajuster les valeurs des pixels en fonction de la nouvelle valeur maximale
	for i := range pgm.data {
		for k := range pgm.data[i] {
			// Calculer la nouvelle valeur ajustée en fonction de la nouvelle valeur maximale
			pgm.data[i][k] = uint8(float64(pgm.data[i][k]) * float64(maxValue) / float64(oldMax))
		}
	}
}

// ToPBM converts the PGM image to PBM.
func (pgm *pgm) ToPBM() *PBM {
	var newNumber string
	if pgm.magicNumber == "P2" {
		newNumber = "P1"
	} else if pgm.magicNumber == "P5" {
		newNumber = "P4"
	}
	NumRows := pgm.width
	NumColumns := pgm.height
	var newData = make([][]bool, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newData[i] = make([]bool, NumRows)
		for j := 0; j < NumRows; j++ {
			// Convertir la valeur de niveau de gris en noir ou blanc
			newData[i][j] = (pgm.data[i][j] < pgm.max/2)
		}
	}
	return &PBM{data: newData, width: NumRows, height: NumColumns, magicNumber: newNumber}
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *pgm) Rotate90CW() {
	NumRows := pgm.width
	NumColumns := pgm.height
	var newData [][]uint8
	for i := 0; i < NumRows; i++ {
		newData = append(newData, make([]uint8, NumColumns))
	}

	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColumns; j++ {
			newData[i][j] = pgm.data[NumRows-j-1][i]
		}
	}
	pgm.data = newData
}
