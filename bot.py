import requests
import os

from dataclasses import dataclass

# We check if the discord token exists 
# to see if  we're just testing or actually
# running the bot
token = os.getenv("DISCORD_TOKEN")

if token:
    environment = "production"
    # Make sure token is treated as a string
    token = str(token)
else:
    environment = "testing"
    print("DISCORD_TOKEN not set")

print("Current environment is: {}".format(environment))

if environment == "production":
    import discord
    import dotenv

    dotenv.load_dotenv()

    bot = discord.Bot()

@dataclass
class Command:
    """
    Struct that holds the information about a command.
    """
    name: str
    desc: str
    resp: str
    url: str
    headers: str
    data: str
    args: int
    user_arg: bool
    arg_names: dict
    arg_defaults: dict
    group: str

@dataclass
class Group:
    """
    Struct that holds the information about a group.
    """
    name: str
    desc: str
    parent: str

def cmd_builder(cmd):
    """
    This function generates the actual function that will be executed at runtime.
    It is not what most people would consider "clean code" but its literally
    generating Python code on the fly so forgive the abysmal syntax.

    I may refactor it if I continue to add functionality.
    Till then good luck.
    """

    cmd_string = "@" + cmd.group + ".command(description='" + cmd.desc + """')
async def """ + cmd.name + """(ctx"""
    for i in range(cmd.args):
        cmd_string += ", " + (cmd.arg_names[i] if i in cmd.arg_names else str("arg" + str(i))) + (str("='" + cmd.arg_defaults[i] + "'") if i in cmd.arg_defaults else "")

    if cmd.args > 0:
        cmd_string += """):
    data = '""" + cmd.data + "'.format("
        for i in range(cmd.args):
            cmd_string += (cmd.arg_names[i] if i in cmd.arg_names else str("arg" + str(i))) + ", "
        if cmd.user_arg:
            cmd_string += "user=ctx.author.name)" 
        else:
            cmd_string = cmd_string[:-2] + ")" 
    elif cmd.args == 0 and cmd.user_arg:
        cmd_string += """):
    data = '""" + cmd.data + "'.format(user=ctx.author.name)"
    else:
        cmd_string += """):
    data = '""" + cmd.data + "'"

    cmd_string += """
    req = requests.post('""" + cmd.url + """', headers=""" + cmd.headers + """, data=data)
    await ctx.respond(f\"""" + cmd.resp + "\")"

    return cmd_string

def group_builder(group):
    """
    This function creates a new command group at runtime if needed.
    I guess I would call this clean although cmd_builder also used to be like this.
    """
    return group.name + " = " + group.parent + ".create_group('" + group.name + "', '" + group.desc + "')"

def get_commands():
    """
    This function parses the environment variables that contain the commands.
    I am honestly okay with this function.  If I add much more to it, it might
    be getting to be too much, but in its current state it seems okay.
    """
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
            args=0,
            user_arg=False,
            arg_names={},
            arg_defaults={},
            group="bot"
        )
        if os.getenv(str(current_prefix + "_HEADERS")):
            cmd.headers = str(os.getenv(str(current_prefix + "_HEADERS")))
        if os.getenv(str(current_prefix + "_DATA")):
            cmd.data = str(os.getenv(str(current_prefix + "_DATA")))
            if os.getenv(str(current_prefix + "_USER_ARG")):
                cmd.user_arg = True
            if os.getenv(str(current_prefix + "_ARGS")):
                cmd.args = int(os.getenv(str(current_prefix + "_ARGS")))
                for x in range(cmd.args):
                    if os.getenv(str(current_prefix + "_ARG_" + str(x) + "_NAME")):
                        cmd.arg_names[x] = str(os.getenv(str(current_prefix + "_ARG_" + str(x) + "_NAME")))
                    if os.getenv(str(current_prefix + "_ARG_" + str(x) + "_DEFAULT")):
                        cmd.arg_defaults[x] = str(os.getenv(str(current_prefix + "_ARG_" + str(x) + "_DEFAULT")))
        if os.getenv(str(current_prefix + "_GROUP")):
            cmd.group = str(os.getenv(str(current_prefix + "_GROUP")))
        commands.append(cmd) 
        i += 1
    return commands

def get_groups():
    """
    This function parses the environment variables that contain the command groups.
    Its pretty similar to get_commands but for groups instead.
    """
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
        if environment == "production":
            exec(group_builder(group))
        else:
            print(group_builder(group))

    for cmd in commands:
        print("Loading command: '" + cmd.name + "' with desc: '" + cmd.desc + "'")
        if environment == "production":
            exec(cmd_builder(cmd))
        else:
            print(cmd_builder(cmd))
    
    print("Success.  Starting dishook...")
    if token:
        bot.run(token)
    
