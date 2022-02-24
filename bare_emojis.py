with open("emojis.txt", 'r') as input:
    with open("bare_emojis.txt", 'w') as output:
        for inputLine in input:
            output.write(inputLine.split("# ")[1].split()[0] + '\n')
