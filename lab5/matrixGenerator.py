from random import randint

def createMatrix(sizeR, sizeC):
    #sizeR, sizeC = var.split(",")
    matrix = []
    for _ in range(int(sizeR)):
        col = []
        for _ in range(int(sizeC)):
            col.append(randint(0, 100))
        matrix.append(col)
    return matrix

def writeMatrix(matrix, sizeR, sizeC, fileName):
    f = open(fileName, "w")
    f.write(sizeR + " " + sizeC + "\n")
    f.write(str(matrix))
    f.close()


def readInput():
	global sizeR0, sizeC0, sizeR1, sizeC1
	
	var0 = raw_input("Please enter the size for the first matrix --> #Rows, #Columns: ")
	var1 = raw_input("Please enter the size for the second matrix --> #Rows, #Columns: ")

	sizeR0, sizeC0 = var0.split(",")
	sizeR1, sizeC1 = var1.split(",")
	if int(sizeC0) != int(sizeR1):
		print "The sizes of Matrix 1: #Columns and Matrix 2: #Rows has to be the same"
		readInput()

sizeR0, sizeC0, sizeR1, sizeC0 = "", "", "", ""
readInput()

matrix0 = createMatrix(sizeR0, sizeC0)
writeMatrix(matrix0, sizeR0, sizeC0, "matrix0.txt")

matrix1 = createMatrix(sizeR1, sizeC1)
writeMatrix(matrix1, sizeR1, sizeC1, "matrix1.txt")
