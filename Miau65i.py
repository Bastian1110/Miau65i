import sys 

ADDRESS_COUNT = 0x02
RESERV_WORDS = ["if","else","while"]

def checkAddresCount():
    global ADDRESS_COUNT
    if(ADDRESS_COUNT == 254):
        ADDRESS_COUNT = 512
    else:
        ADDRESS_COUNT+=2

class Equal:
    def getTable(scope,arguments, optype):
        out = {"SCOPE" : scope, "ADDRESS" : str(hex(ADDRESS_COUNT))[2:], "VALUE" : arguments, "OPTYPE" : optype}
        checkAddresCount()
        return out

    def getAsm(table,data,name):
        if(data["OPTYPE"] == "DECL"):
            if(data["VALUE"][0].isdigit()):
                msb = "{:04x}".format(int(data["VALUE"][0]))[:2]
                lsb = "{:04x}".format(int(data["VALUE"][0]))[2:]
                print("  lda #$" + msb)
                print("  sta $" + data["ADDRESS"] + "      ; " + name)
                print("  lda #$" + lsb)
                print("  sta $" + data["ADDRESS"] + " + 1\n")
            else:
                print("  lda $" + table[data["VALUE"][0]]["ADDRESS"])
                print("  sta $" + data["ADDRESS"] + "      ; " + name)
                print("  lda $" + table[data["VALUE"][0]]["ADDRESS"] + " + 1")
                print("  sta $" + data["ADDRESS"] + " + 1\n")
        elif(data["OPTYPE"] == "MATH"):
            if('+' in data["VALUE"]):
                Plus.getAsm(table,data["VALUE"][0],data["VALUE"][2])
                print("  sty $" + data["ADDRESS"] + "      ; " + name)
                print("  sta $" + data["ADDRESS"] + " + 1\n")

class Plus:
    def getTable(scope,argumentOne,argumentTwo):
        out = {"SCOPE" : scope, "ADDRESS" : str(hex(ADDRESS_COUNT)), "VALUEONE" : argumentOne,"VALUETWO": argumentTwo}
        return out

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

            
class Iff:
    def getTable(statement,beginIf,endIf,beginElse,endElse):
        out = {"IFBEGIN":beginIf,"IFEND":endIf,"ELSEBEGIN":beginElse,"ELSEEND":endElse,"STATEMENT":statement}
        return out

    def searchInfo(program):
        out = [0,0,0,0] #IfBegin , IfEnd , ElseBegin, ElseEnd
        stack = []
        res = []
        for l in range(len(program)):
            if('{' in program[l]):
                stack.append(program[l])

            elif('}' in program[l]):
                res.append([stack.pop(),l])
        print(stack)
        print(res)
        stack = []
        res = []
        print("---------------")
        for l in range(len(program)):
            if('if' in program[l]):
                stack.append(program[l])

            elif('else' in program[l]):
                res.append([stack.pop(),l])
        print(stack)
        print(res)
        return out
                



            

def createStateTable(code):
    x = 0
    table = {}
    Iff.searchInfo(code)
    return table
    
    



def getOperand(line):
    operands = ['=','+','-',"p"]
    out = {}
    for i in operands:
        out[i] = []
        x = 0
        for w in line:
            if(w == i):
                out[i].append(x)
            x+=1
    return out


def createSymbolTable(code):
    table = {}
    for l in code:
        operands = getOperand(l)
        if(len(operands['=']) != 0):
            if(len(operands['+']) != 0 or len(operands['-']) != 0):
                table[l[0]] = Equal.getTable(1,l[2:],"MATH")
            else:
                table[l[0]] = Equal.getTable(1,l[2:],"DECL")
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
            program[x] = program[x].split(' ')
            while(' ' in program[x]):
                blank = program[x].index(' ')
                program[x].pop(blank)
            while('' in program[x]):
                blank = program[x].index('')
                program[x].pop(blank)
        else:
            program[x] = handleKeyWords(program[x])
    return program

def parseCode(code,table):
    for l in code:
        pass
        



def assemblyCode():
    program = readFile("foo.miau")
    symTable = createSymbolTable(program)
    stateTable = createStateTable(program)

    parseCode(program,symTable)
#    for i in symTable:
#        print(i,symTable[i])




if __name__ == "__main__":
    #TODO #1  Create a proper command line for the compiler
    assemblyCode()



