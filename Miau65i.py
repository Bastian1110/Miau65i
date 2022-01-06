import sys 

ADDRESS_COUNT = 0x0
def checkAddresCount():
    global ADDRESS_COUNT
    if(ADDRESS_COUNT == 254):
        ADDRESS_COUNT = 512
    else:
        ADDRESS_COUNT+=1

class Equal:
    def getTable(scope,arguments):
        out = {"SCOPE" : scope, "ADDRESS" : str(hex(ADDRESS_COUNT)), "VALUE" : arguments}
        checkAddresCount()
        return out

class Plus:
    def getTable(scope,argumentOne,argumentTwo):
        out = {"SCOPE" : scope, "ADDRESS" : str(hex(ADDRESS_COUNT)), "VALUEONE" : argumentOne,"VALUETWO": argumentTwo}
        return out

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
        if(operands['='] != 0):
            pass #TODO



def readFile(filePath):
    code = open(filePath,'r')
    program = code.read()
    program = program.split('\n')
    code.close()
    for x in range(len(program)):
        program[x] = program[x].split(' ')
        while(' ' in program[x]):
            blank = program[x].index(' ')
            program[x].pop(blank)
    return program




if __name__ == "__main__":
    program = readFile("test.miau")
    createSymbolTable(program)




