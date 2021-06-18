# Introduction
'mango' is a CLI based task parser and tester for popular online judge 'Codeforces'. It supports the regular contest and gym.

# Download

# Installation
Windows:<br>
1. Keep 'mango.exe' in any folder you prefer
2. Add the folder path from step 1 to System Variable Path. (How to add path in System Var? See here: https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/)

Mac:<br>
1. Keep 'mango' executable file at /usr/local/bin folder
2. Change the permission of the file by 'chmod +x mango' command.

Linux:<br>
May be similar to mac. Didn't try in linux.

# Configuration
1. Set default programs to open .cpp & .json files. (How to change default programs in windows? See here: https://www.digitaltrends.com/computing/how-to-change-file-associations/)
2. Open cmd prompt from anywhere and run 'mango configure'. It will open config.json file. Or go to AppData>Roaming>mango, you'll find the config.json file. Now configure as you prefer. But DO NOT CHANGE the OJ and Host property.
![](./assets/Configuration_Json.png)
3. Set 'Workspace' as the full path of the folder where all of the contest sources and testcases will be stored
4. Set 'TemplatePath' as the full path of your template file. If you ommit TemplatePath, a default template will be created for the source file.
5. Set 'Author' as your username/handle or anything name you prefer. It will be used in your template.
6. Enjoy!



# Workspace Structure:
```    
workspace
├── codeforces
│   └── 1521
│       ├── src
│       │   ├── a.cpp
│       │   ├── b.cpp
│       │   │
│       │   │
│       │   └── e.cpp
│       └── testcase
│       │   ├── a.json
│       │   ├── b.json
│       │   │
│       │   │
│       │   └── e.json
```

# Command Format
`mango <command> <argument>`

=> Only `configure` and `version` commands don't need any argument.<br>
=> For other comamnds, the argument format is `<contest_id><problem_id>`. But both `<contest_id>`and `<problem_id>` are optional for corresponding command.


# Available Commands as example
1. `mango setc 1521`: sets current working contest ID
2. `mango configure`: opens the config.json file to update & save configuration

3. `mango parse 1521`: command 1 + parses samples of all the problems for specified contest ID
4. `mango parse 1521A`: command 1 + parses samples of Problem A for specified contest ID

5. `mango source 1521A`: creates source file of Problem A for specified contest ID
6. `mango source A`: creates source file of Problem A for current working contest ID

7. `mango open 1521`: opens all the source files in the default editor for specified contest ID
8. `mango open 1521A`: opens source file of Problem A in the default editor for specified contest ID
9. `mango open A`: opens source file of Problem A in the default editor for current working contest ID

10. `mango create 1521`: combination of commands (1, 3, 5, 7) for specified contest ID
11. `mango create 1521A`: combination of commands (1, 4, 6, 8) for Problem A

12. `mango compile 1521A`: compiles source file of Problem A for specified contest ID
13. `mango compile A`: compiles source file of Problem A for current working contest ID

14. `mango test 1521A`: command 12 + tests Problem A for specified contest ID
15. `mango test A`: command 13 + tests Problem A for current working contest ID

16. `mango version`: prints the current mango version


# Configure mango with Sublime Text
We can directly test our code from Sublime Text through the custom build system. To configure the sublime build system in Windows, please follow the instructions below..
1. From Sublime Text menubar, go to Tools > Build System > New Build System
2. Copy & Paste the following code
  ```
  {
	"shell_cmd": "mango test \"${file_base_name}\"",
	"file_regex": "^(..[^:]*):([0-9]+):?([0-9]+)?:? (.*)$",
	"working_dir": "${file_path}",
	"selector": "source.c, source.cc, source.c++, source.cpp",

	"variants":
	[
		{
			"name": "Interactive",
			"shell_cmd": "start cmd /k mango test \"${file_base_name}\" & pause & exit\""
		}
	]
}
  ```
3. Press Ctrl+S to save the file. Give it a name like `cpp_mango_tester.sublime-build` and save.
4. Go to Tools > Build System. You will find your build name(i.e `cpp_mango_tester`) in the list. Select the build name.
5. Now create a problem or contest by 'mango create' command and write your code.
6. This step needs to be done only once. Goto Tools > Build With (Shortcut Ctrl+Shift+B), a pop-up will appear on top-center of the screen. Select `cpp_mango_tester - Interactive` from the pop-up.
7. To test your code from Sublime Text, go to Tools > Build.
