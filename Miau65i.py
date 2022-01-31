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
                print("  sty $" + table[name]["ADDRESS"] + "      ; " + name)
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
    def getTable(statement,begin,end,elsee):
        return {"IFBEGIN":begin,"IFEND":end,"STATEMENT":statement,"ELSEKEY" : elsee}
    def getAsm(code,tableSym,tableFlow,labelCounter,elsee):
        statement = Boolean.getAsm(tableSym,tableFlow["IF"+str(labelCounter)]["STATEMENT"])
        if(elsee):
            if(statement[1] == '>'):
                print("  bpl IF"+str(labelCounter))
                print("")
        else:
            if(statement[1] == '>'):
                print("  bpl END"+str(labelCounter))


class Elsee:
    def getTable(iff,begin,end):
        return {"ELBEGIN":begin,"ELEND" : end,"IFKEY" : iff}
    def getAsm(code,table):
        Boolean.getAsm()

class Whilee:
    def getTable(statement,begin,end):
        return {"WHBEGIN":begin,"WHEND":end,"STATEMENT":statement}




            


def getBlock(blocks,index):
    for i in blocks:
        if(i[0] == index):
            return i



def generateThreeAddresCode(code):
    tac = []
    ready = []
    labelCounter = 0
    for i in range(len(code)):
        print(i,code[i])

    stack = []
    blocksLimits = []
    for l in range(len(code)):
        if('{' in code[l]):
            stack.append(l)
        elif('}' in code[l]):
            blocksLimits.append([stack.pop(),l])

    print(blocksLimits)
    print("---------")
    
    
    for i in range(len(code)):
        if(i not in ready):
            if("if" in code[i]):
                ifBlock = getBlock(blocksLimits,i)
                if("else" in code[ifBlock[1]+1]):

                    tac.append(code[ifBlock[0]][:2] + ["go to","L"+str(labelCounter)+':'])
                    elseBlock = getBlock(blocksLimits,ifBlock[1]+1)

                    for l in range(elseBlock[0]+1,elseBlock[1]):
                        tac.append(code[l])
                        ready.append(l)
                    tac.append(["go to","L"+str(labelCounter+1)])
                    ready.append(elseBlock[0])
                    ready.append(elseBlock[1])

                    tac.append(["L"+str(labelCounter)+':'])
                    for l in range(ifBlock[0]+1,ifBlock[1]):
                        tac.append(code[l])
                        ready.append(l)
                    ready.append(ifBlock[0])
                    ready.append(ifBlock[1])
                    labelCounter+=1

                    tac.append(["L"+str(labelCounter)+':'])
                    labelCounter+=1
                else:
                    pass
            if i not in ready:
                tac.append(code[i])
    
    for i in tac:
        print(i)


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
    out = []
    for x in range(len(program)):
        if program[x]:
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
            out.append(program[x])
    
    return out

def parseCode(code,symbols):
    pass
        
        



def assemblyCode(path):
    program = readFile(path)    
    symTable = createSymbolTable(program)
    intermidiateCode = generateThreeAddresCode(program)
    parseCode(program,symTable)




if __name__ == "__main__":
    #TODO #2 Create a proper command line for the compiler
    assemblyCode("foo.miau")



