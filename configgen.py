import sys
import random

masterip = sys.argv[1]
masterport = int(sys.argv[2])
myname = sys.argv[3]
myip = sys.argv[4]
portbase = int(sys.argv[5])

genstr = '''#Here it is

Node
    NodeName=master
    NodeType=master
    NodeAddr='''+masterip+'''
    SendPort='''+str(masterport)+'''
    RecvPort='''+str(masterport+1)+'''
    lsp='''+str(masterport+2)+'''
    lcp='''+str(masterport+3)+'''
    dsp='''+str(masterport+4)+'''
    BootstrapPort='''+str(masterport+5)+'''
    Effort=400
    GroupSize=4
End

Node
    NodeName='''+str(myname)+'''
    NodeType=none
    NodeAddr='''+myip+'''
    NodeGroup=group
    SendPort='''+str(portbase)+'''
    RecvPort='''+str(portbase+1)+'''
    lsp='''+str(portbase+2)+'''
    lcp='''+str(portbase+3)+'''
    dsp='''+str(portbase+4)+'''
    BootstrapPort='''+str(portbase+5)+'''
    Effort='''+str(random.randrange(100,1000,50))+'''
End
'''

print genstr
