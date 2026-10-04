# vhmm

> [!CAUTION]
> vhmm is made for me and me only, don't expect it to work out of the box.

Tiny Valheim mod manager for Linux.


vhmm does the following:
* keeps track of a local index of thunderstore from which to download plugins from.
* maintains all your plugins that are locally installed in your `<GAME>/BepInEx/plugins` directory.
They are assumed to extracted to their own directory, so your `plugins` directory should look something like this:
```
../Valheim/BepIndex/plugins/
../Valheim/BepIndex/plugins/denikson-BepInExPack_Valheim-5.4.2350
../Valheim/BepIndex/plugins/shudnal-ConditionalConfigSync-1.0.10
../Valheim/BepIndex/plugins/shudnal-ExtraSlots-1.2.16
../Valheim/BepIndex/plugins/ValheimModding-YamlDotNet-16.3.1
```
* resolves any dependencies for any plugin you've installed and attempts to
  download it and update it whenever a new version has been released.
* deletes old versions of any installed plugin. Do keep in mind that
  dependencies that are no longer used by a plugin are not removed
  automatically.


## Usage

> [!WARNING]
> vhmm does not manage BepInEx, so you are expected to have installed this yourself.

Run `vhmm` with the `-h` flag which yields the following:
```
Usage of vhmm:
  -dry-run
    	Only display changes
  -game string
    	Path to the root game directory
  -repository string
    	URL of the repository (default "https://thunderstore.io/c/valheim/api/v1/package/")
  -script string
    	Optional; if supplied and vhmm finishes without any issues, run script after
```

> [!NOTE]
> The only flag that is required for vhmm to work is the `-game` flag which should point to your *root* directory of Valheim.


### Example

1. Download any plugin of your choice from thunderstore.io (for example `ExtraSlots`)
1. Extract all of the contents from the zip file to its own directory (for example `shudnal-ExtraSlots-1.2.16`)
1. Move the entire plugin directory to `</path/to/Valheim>/BepInEx/plugins`
1. Copy the `vhmm` executable in to `</path/to/Valheim>`
1. Run `vhmm -game </path/to/Valheim> -script </path/to/Valheim>/start_game_bepinex.sh` to tell vhmm where Valheim is installed and optionally launch Valheim with BepInEx enabled after all mods are confirmed to be up-to-date.

Having to provide flags everytime you launch vhmm is awkward and
annoying so you can create a `config.json` file in the same directory as
vhmm is launched where you can provide default values for both of the flags
mentioned above, like this:
```
	{
		"game_dir": "/path/to/Valheim",
		"script": "/path/to/Valheim/start_game_bepinex.sh"
	}
```
