from random import randint

def createMatrix(var):
    sizeR, sizeC = var.split(",")
    matrix = []
    for _ in range(int(sizeR)):
        col = []
        for _ in range(int(sizeC)):
            col.append(randint(0, 100))
        matrix.append(col)
    return sizeR, sizeC, matrix

def writeMatrix(matrix, sizeR, sizeC, fileName):
    f = open(fileName, "w")
    f.write(sizeR + " " + sizeC + "\n")
    f.write(str(matrix))
    f.close()

var0 = input("Please enter the size for the first matrix --> #Rows, #Columns: ")
var1 = input("Please enter the size for the second matrix --> #Rows, #Columns: ")
# Si les mides no quadren, treure un missatge per pantalla demanant de nou les mides de la segona matriu

sizeR0, sizeC0, matrix0 = createMatrix(var0)
writeMatrix(matrix0, sizeR0, sizeC0, "matrix0.txt")

sizeR1, sizeC1, matrix1 = createMatrix(var1)
writeMatrix(matrix1, sizeR1, sizeC1, "matrix1.txt")