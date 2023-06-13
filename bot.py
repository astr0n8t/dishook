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
    group: str

@dataclass
class Group:
    name: str
    desc: str
    parent: str

def cmd_builder(cmd):

    cmd_string = "@" + cmd.group + ".command(description='" + cmd.desc + """')
async def """ + cmd.name + """(ctx"""

    cmd_string += """):
    req = requests.post('""" + cmd.url + """', headers=""" + cmd.headers + """, data=""" + cmd.data + """)
    await ctx.respond(f\"""" + cmd.resp + "\")"

    return cmd_string

def group_builder(group):
    return group.name + " = " + group.parent + ".create_group('" + group.name + "', '" + group.desc + "')"

def get_commands():
    commands = []
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
            data="{}",
            group="bot"
        )
        if os.getenv(str(current_prefix + "_HEADERS")):
            cmd.headers = str(os.getenv(str(current_prefix + "_HEADERS")))
        if os.getenv(str(current_prefix + "_DATA")):
            cmd.headers = str(os.getenv(str(current_prefix + "_DATA")))
        if os.getenv(str(current_prefix + "_GROUP")):
            cmd.group = str(os.getenv(str(current_prefix + "_GROUP")))
        commands.append(cmd) 
        i += 1
    return commands

def get_groups():
    groups = []
    dishook_prefix = "DISHOOK_GROUP_"
    i = 0
    while os.getenv(str(dishook_prefix + str(i))):
        current_prefix = str(dishook_prefix + str(i))
        group = Group(
            name=str(os.getenv(current_prefix)),
            desc=str(os.getenv(str(current_prefix + "_DESCRIPTION"))),
            parent="bot"
        )
        if os.getenv(str(current_prefix + "_PARENT")):
            group.parent = str(os.getenv(str(current_prefix + "_PARENT")))
        groups.append(group)
        i += 1
    return groups

if __name__ == '__main__':

    print("Beginning to look for commands designated by DISHOOK_COMMAND_#")

    commands = get_commands()
    groups = get_groups()

    for group in groups:
        print("Loading groups: '" + group.name + "' with desc: '" + group.desc + "'")
        exec(group_builder(group))

    for cmd in commands:
        print("Loading command: '" + cmd.name + "' with desc: '" + cmd.desc + "'")
        exec(cmd_builder(cmd))
    
    print("Success.  Starting dishook...")
    bot.run(token)
    
