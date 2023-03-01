import requests
import discord
import dotenv
import os

from dataclasses import dataclass

dotenv.load_dotenv()
token = str(os.getenv("DISCORD_TOKEN"))

bot = discord.Bot()

@dataclass
class Command:
    name: str
    desc: str
    resp: str
    url: str

def cmd_builder(name, desc, resp, url):

    return"@bot.command(description='" + desc + """')
async def """ + name + """(ctx):
    req = requests.post('""" + url + """')
    await ctx.respond(f\"""" + resp + "\")"

def get_commands():
    command_list = []
    dishook_prefix = "DISHOOK_COMMAND_"
    i = 0
    while os.getenv(str(dishook_prefix + str(i))):
        cmd = Command(
            name=str(os.getenv(str(dishook_prefix + str(i)))),
            desc=str(os.getenv(str(dishook_prefix + str(i) + "_DESCRIPTION"))),
            resp=str(os.getenv(str(dishook_prefix + str(i) + "_RESPONSE"))),
            url=str(os.getenv(str(dishook_prefix + str(i) + "_URL")))
        )
        command_list.append(cmd) 
        i += 1
    return command_list

if __name__ == '__main__':

    print("Beginning to look for commands designated by DISHOOK_COMMAND_#")

    commands = get_commands()

    for cmd in commands:
        print("Loading command: '" + cmd.name + "' with desc: '" + cmd.desc + "'")
        exec(cmd_builder(cmd.name, cmd.desc, cmd.resp, cmd.url))
    
    print("Success.  Starting dishook...")
    bot.run(token)
    
