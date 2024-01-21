package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func bitsToString(byteSlice []byte) string {
	var result string
	for _, b := range byteSlice {
		for i := 7; i >= 0; i-- {
			bitValue := (b & (1 << uint(i))) >> uint(i)
			result += strconv.Itoa(int(bitValue))
		}
	}
	return result
}

func stringToBool(line string, width int, height int) []bool {
	caracteres := strings.Fields(line)
	if len(caracteres) != width {
		fmt.Println("Erreur dans la ligne récupérée")
		return nil
	}

	// Initialiser la tranche de données avec la bonne taille
	data := make([]bool, width)

	// Utiliser une boucle for pour parcourir les caractères et remplir la tranche de données
	for i := 0; i < width; i++ {
		if caracteres[i] == "0" {
			data[i] = false
		} else if caracteres[i] == "1" {
			data[i] = true
		} else {
			fmt.Println("Caractère non valide trouvé dans la ligne")
			return nil
		}
	}

	return data
}

func ReadPBM(filename string) (*PBM, error) {

	var pbmR PBM
	var pbm *PBM = &pbmR
	var verifMN bool = false
	var verifDim bool = false

	file, err := os.Open(filename)
	if err != nil {
		// Gérer l'erreur s'il y en a une
		fmt.Println("Erreur lors de l'ouverture du fichier:", err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// Si la ligne est vide, cela peut signifier que la fin du fichier est atteinte
			break
		}

		if !verifMN {
			pbm.magicNumber = line
			verifMN = true
			continue
		}
		fmt.Println("Line:", line)

		if line[0] == '#' {
			continue
		}

		if !verifDim {
			curLine := strings.Fields(line)
			pbm.height, err = strconv.Atoi(curLine[0])
			if err != nil {
				fmt.Println("erreur dans la lecture de la longueur")
			}
			pbm.width, err = strconv.Atoi(curLine[1])
			if err != nil {
				fmt.Println("erreur dans la lecture de la largeur")
			}
			pbm.data = make([][]bool, pbm.height)
			for k := range pbm.data {
				pbm.data[k] = make([]bool, pbm.width)
			}
			verifDim = true
			continue
		}

		if pbm.magicNumber == "P1" {
			if i < pbm.height {
				pbm.data[i] = stringToBool(line, 15, 15)
				i++
			}
		}

		if pbm.magicNumber == "P4" {
			if i >= pbm.height {
				// Atteint la fin du tableau, peut-être sortir de la boucle
				break
			}
			fmt.Println([]rune(line))
			lineUpdate := []byte(line)
			fmt.Println(lineUpdate)
			stringBytes := bitsToString(lineUpdate)
			fmt.Println(stringBytes)
			data := make([][]bool, pbm.height)
			compteur := 0
			for i := range data {
				data[i] = make([]bool, pbm.width)
				for j := range data[i] {
					if stringBytes[compteur] == 1 {
						data[i][j] = true
					} else if stringBytes[compteur] == 0 {
						data[i][j] = false
					}
					compteur += 1
				}
			}
			pbm.data = data
		}
	}
	return pbm, err
}

func (pbm *PBM) Size() (int, int) {
	return pbm.height, pbm.width
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erreur lors de la création/ouverture du fichier :", err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(pbm.magicNumber + "\n")
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du nombre magique dans le fichier :", err)
		return err
	}

	_, err = writer.WriteString(strconv.Itoa(pbm.height) + " " + strconv.Itoa(pbm.width) + "\n")
	if err != nil {
		fmt.Println("Erreur lors de l'écriture de la hauteur et de la largeur dans le fichier :", err)
		return err
	}

	for i := 0; i < pbm.height; i++ {
		sliceUpdate := make([]string, 0, pbm.width)
		for k := 0; k < pbm.width; k++ {
			if pbm.data[i][k] {
				sliceUpdate = append(sliceUpdate, "0 ")
			} else {
				sliceUpdate = append(sliceUpdate, "1 ")
			}
		}
		stringUpdate := strings.Join(sliceUpdate, " ")
		_, err = writer.WriteString(stringUpdate + "\n")
		if err != nil {
			fmt.Println("Erreur lors de l'écriture des données dans le fichier :", err)
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Erreur lors du vidage du buffer dans le fichier :", err)
		return err
	}

	return nil
}

func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for k := 0; k < pbm.width; k++ {
			if pbm.data[i][k] {
				pbm.data[i][k] = false
			} else {
				pbm.data[i][k] = true
			}
		}
	}
}

func (pbm *PBM) Flop() {
	if pbm.width <= 0 || pbm.height <= 0 {
		fmt.Println("Invalid dimensions")
		return
	}

	sliceUpdate := make([][]bool, pbm.height)
	for i := range sliceUpdate {
		sliceUpdate[i] = make([]bool, pbm.width)
	}

	for i := 1; i <= pbm.height; i++ {
		for k := 0; k < pbm.width; k++ {
			sliceUpdate[i-1][k] = pbm.data[pbm.height-i][k]
		}
	}
	pbm.data = sliceUpdate
}

func (pbm *PBM) Flip() {

	fmt.Println(pbm.data)
	if pbm.width <= 0 || pbm.height <= 0 {
		fmt.Println("Invalid dimensions")
		return
	}

	sliceUpdate := make([][]bool, pbm.height)
	for i := range sliceUpdate {
		sliceUpdate[i] = make([]bool, pbm.width)
	}

	for i := 0; i < pbm.height; i++ {
		for k := 0; k < pbm.width; k++ {
			sliceUpdate[i][k] = pbm.data[i][pbm.width-k-1]
		}
	}
	pbm.data = sliceUpdate

}

func (pbm *PBM) SetMagicNumber(magicNumber string) {

	if !(magicNumber == "P1" || magicNumber == "P4") {
		fmt.Println("Mauvaise saisie pour magicNumber")
	} else {
		pbm.magicNumber = magicNumber
		fmt.Println("Succed")
	}

}
