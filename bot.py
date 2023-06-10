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
    headers: str
    data: str

def cmd_builder(name, desc, resp, url, headers, data):

    return"@bot.command(description='" + desc + """')
async def """ + name + """(ctx):
    req = requests.post('""" + url + """', headers=""" + headers + """, data=""" + data + """)
    await ctx.respond(f\"""" + resp + "\")"

def get_commands():
    command_list = []
    dishook_prefix = "DISHOOK_COMMAND_"
    i = 0
    while os.getenv(str(dishook_prefix + str(i))):
        current_prefix = str(dishook_prefix + str(i))
        cmd = Command(
            name=str(os.getenv(current_prefix)),
            desc=str(os.getenv(str(current_prefix + "_DESCRIPTION"))),
            resp=str(os.getenv(str(current_prefix + "_RESPONSE"))),
            url=str(os.getenv(str(current_prefix  + "_URL"))),
            headers="{}",
            data="{}"
        )
        if os.getenv(str(dishook_prefix + "_HEADERS")):
            cmd.headers = str(os.getenv(str(dishook_prefix + "_HEADERS")))
        if os.getenv(str(dishook_prefix + "_DATA")):
            cmd.headers = str(os.getenv(str(dishook_prefix + "_DATA")))
        command_list.append(cmd) 
        i += 1
    return command_list

if __name__ == '__main__':

    print("Beginning to look for commands designated by DISHOOK_COMMAND_#")

    commands = get_commands()

    for cmd in commands:
        print("Loading command: '" + cmd.name + "' with desc: '" + cmd.desc + "'")
        exec(cmd_builder(cmd.name, cmd.desc, cmd.resp, cmd.url, cmd.headers, cmd.data))
    
    print("Success.  Starting dishook...")
    bot.run(token)
    
