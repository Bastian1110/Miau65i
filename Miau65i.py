import sys 
import re

ADDRESS_COUNT = 0x02
RESERV_WORDS = ["if","else","while"]
OPERATIONS = ['+','-']
BOOLEAN_OP = ['<','>','==','!=']
TYPES = ["int"]

def checkAddresCount(sum):
    global ADDRESS_COUNT
    if(ADDRESS_COUNT == 254):
        ADDRESS_COUNT = 512
    else:
        ADDRESS_COUNT+=sum


class Allocation:
    def getTable(type,scope):
        out = {"ADDRESS" :  str(hex(ADDRESS_COUNT))[2:],"TYPE" : type, "SCOPE" : scope}
        checkAddresCount(2)
        return out


class Equal:
    global OPERATIONS
    def getAsm(table,data):
        if(data[1] in table):
            value = data[3:]
            name = data[1]
        else:
            value = data[2:]
            name = data[0]
        if(len(value) == 1):
            if(value[0].isdigit()):
                msb = "{:04x}".format(int(value[0]))[:2]
                lsb = "{:04x}".format(int(value[0]))[2:]
                print("  lda #$" + msb)
                print("  sta $" + table[name]["ADDRESS"] + "      ; " + name)
                print("  lda #$" + lsb)
                print("  sta $" + table[name]["ADDRESS"] + " + 1\n")
            else:
                print("  lda $" + table[value[0]]["ADDRESS"])
                print("  sta $" + table[name]["ADDRESS"] + "      ; " + name)
                print("  lda $" + table[value[0]]["ADDRESS"] + " + 1")
                print("  sta $" + table[name]["ADDRESS"] + " + 1\n")
        else:
            if('+' in value):
                Plus.getAsm(table,value[0],value[2])
                print("  sty $" + table[name]["ADDRESS"] + "      ; " + data[1])
                print("  sta $" + table[name]["ADDRESS"]+ " + 1\n")

class Plus:
    def getAsm(table,argumentOne,argumentTwo):
        typeP = [argumentOne.isdigit(),argumentTwo.isdigit()]
        if (not typeP[0] and not typeP[1]):
            print("  lda $" + table[argumentOne]["ADDRESS"])
            print("  adc $" + table[argumentTwo]["ADDRESS"])
            print("  tay         ; " + argumentOne + " + " + argumentTwo)
            print("  lda $" + table[argumentOne]["ADDRESS"] + " + 1")
            print("  adc $" + table[argumentTwo]["ADDRESS"] + " + 1")
        else:
            direct = argumentTwo if typeP[1] else argumentOne
            indirect = argumentTwo if typeP[0] else argumentOne
            msb = "{:04x}".format(int(direct))[:2]
            lsb = "{:04x}".format(int(direct))[2:]
            print("  lda #$" + msb)
            print("  adc $" + table[indirect]["ADDRESS"])
            print("  tay         ; " + argumentOne + " + " + argumentTwo)
            print("  lda #$" + lsb)
            print("  adc $" + table[indirect]["ADDRESS"] + " + 1")

class Boolean:
    def getAsm(symbolsTable,statement):
        global BOOLEAN_OP
        statement = re.split('([^a-zA-Z0-9])',statement)
        while ' ' in statement:
            blank = statement.index(' ')
            statement.pop(blank)
        while '' in statement:
            blank = statement.index('')
            statement.pop(blank)
        if(statement.count('=') == 2):
            statement.pop(statement.index('='))
            statement[1] = '=='
        operator = statement[1]

        if(operator == '<'):
            print("Menor Que")
        elif(operator == '>'):
            print("  sec")
            print("  lda $" + symbolsTable[statement[0]]["ADDRESS"])
            print("  sbc $" + symbolsTable[statement[2]]["ADDRESS"])
            print("  tay         ; " + statement[0] + " < " + statement[2])
            print("  lda $" + symbolsTable[statement[0]]["ADDRESS"] + " + 1")
            print("  sbc $" + symbolsTable[statement[2]]["ADDRESS"] + " + 1")


        elif(operator == '=='):
            print("IGUAL")

        return statement
       

            



            
class Iff:
    def getTable(statement,begin,end):
        return {"IFBEGIN":begin,"IFEND":end,"STATEMENT":statement}
    def getAsm(code,tableSym,tableFlow,labelCounter):
        statement = Boolean.getAsm(tableSym,tableFlow["IF"+str(labelCounter)]["STATEMENT"])
        if(statement[1] == '>'):
            print("  bpl IF"+str(labelCounter))
            print("")
            print("ELSE0:")


class Elsee:
    def getTable(iff,begin,end):
        return {"ELBEGIN":begin,"ELEND" : end,"IFKEY" : iff}
    def getAsm(code,table):
        Boolean.getAsm()

class Whilee:
    def getTable(statement,begin,end):
        return {"WHBEGIN":begin,"WHEND":end,"STATEMENT":statement}




            

def createStateTable(code):
    table = {}
    stack = []
    blocksLimits = []
    label = 0
    ifCounter = []
    for l in range(len(code)):
        if('{' in code[l]):
            stack.append(l)
        elif('}' in code[l]):
            blocksLimits.append([stack.pop(),l])
    blocksLimits = sorted(blocksLimits, key=lambda x: x[0])
    for i in blocksLimits:
        if("if" in code[i[0]]):
            table["IF"+str(label)] = Iff.getTable(code[i[0]][1],i[0],i[1])
            if("else" in code[i[1]+1]):
                ifCounter.append(label)
            label+=1
        elif("else" in code[i[0]]):
            table["ELSE"+str(ifCounter[-1])] = Elsee.getTable("IF"+str(ifCounter[-1]),i[0],i[1])
            ifCounter.pop()
        elif("while" in code[i[0]]):
            table["WHILE" + str(label)] = Whilee.getTable(code[i[0]][1],i[0],i[1])

            
    return table
    
    

def createSymbolTable(code):
    global TYPES
    table = {}
    for l in code:
        for ty in TYPES:
            if(ty in l):
                table[l[1]] = Allocation.getTable(ty,1)
    return table


def handleKeyWords(line):
    line = line + '('
    banned = ['(',')']
    out = []
    word = ""
    for i in line:
        if(i not in banned):
            word = word + i
        else:
            out.append(word)
            word = ""
    if("if" in line):
        out[0] = out[0].replace(' ','')
    if("else" in line):
        out = ["else","{"]
    return out


def readFile(filePath):
    global RESERV_WORDS
    code = open(filePath,'r')
    program = code.read()
    program = program.split('\n')
    code.close()
    for x in range(len(program)):
        reservFlag = not(any(check in program[x] for check in RESERV_WORDS))
        if reservFlag:
            program[x] = re.split('([^a-zA-Z0-9])',program[x])
            while(' ' in program[x]):
                blank = program[x].index(' ')
                program[x].pop(blank)
            while('' in program[x]):
                blank = program[x].index('')
                program[x].pop(blank)
        else:
            program[x] = handleKeyWords(program[x])
    return program

def parseCode(code,tableS,tableC):
    labelCounter = 0
    for l in code:
        if('=' in l):
            Equal.getAsm(tableS,l)

        elif("if" in l):
            Iff.getAsm(code,tableS,tableC,labelCounter)
            labelCounter+=1
            
        elif("else" in l):
            print("Aqui va un Else")
        



def assemblyCode():
    program = readFile("foo.miau")
    symTable = createSymbolTable(program)
    stateTable = createStateTable(program)
    parseCode(program,symTable,stateTable)

    for i in stateTable:
        print(i,stateTable[i])




if __name__ == "__main__":
    #TODO #2 Create a proper command line for the compiler
    assemblyCode()



