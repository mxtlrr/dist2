key_dict = {}
class INIParser:
    def __init__(self) -> None:
        pass

    def ParseINIFile(name: str):
        for line in open(name, "r"):
            if line[0] != '[':
                spaced = line.replace("\n","").split(" ")
                key_dict.update({spaced[0]: spaced[2]})

    def GetINIKey(key: str) -> str:
        return key_dict[key]
