from dataclasses import dataclass


@dataclass
class Command:
    id: str 
    cmdText: str
    nextCmd: None
    prevCmd: None

def PrintCommandList(cmdList: list[Command]) -> None:
    for cmd in cmdList:
        print(str(cmd))

def RewindCommandHistory(cmd: Command) -> Command:
    current: Command = cmd
    if current.prevCmd is None:
        return current
    else:
        RewindCommandHistory(current.prevCmd)


def TraverseCommandHistory(cmd: Command) -> list[Command]:
    current:Command = None
    history:list = []

    if (cmd.prevCmd is not None):
        #print("ERR - must be root Command")
        current = RewindCommandHistory(current)
    else:
        current = cmd

    while current is not None:
        history.append(cmd)
        current = cmd.nextCmd
    history.append(cmd) # make sure to at least capture the command passed (the list may just contain a single element)
    
    PrintCommandList(history)

def test():
    # TODO, do some more work with linked lists

    #TraverseCommandHistory(c0)
    pass

if __name__ == '__main__':
    test()
