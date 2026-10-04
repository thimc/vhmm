# vhmm

[!CAUTION]
vhmm is made for me and me only, don't expect it to work out of the box.

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

Running `vhmm` with the `-h` flag yields:

	Usage of vhmm:
	  -dry-run
		Only display changes
	  -game string
		Path to the root game directory
	  -repository string
		URL of the repository (default "https://thunderstore.io/c/valheim/api/v1/package/")
	  -script string
		Optional; if supplied and vhmm finishes without any issues, run script after.

The only flag that is required for vhmm to work is the `-game` flag,
which should point to your *root* directory of Valheim, that is the
directory where there should be another directory called `BepInEx`
and in that `plugins`.

Another useful flag is the `-script` flag which launches any program
you give it when vhmm finishes synchronizing and updating your plugins.
My personal use case for this is to launch Valheim with BepInEx enabled.

Having to provide flags everytime you launch vhmm is awkward and
annoying so you can create a `config.json` file in the same directory as
hmm is launched where you can provide default values for both of the flags
mentioned above, like this:

	{
		"game_dir": "/path/to/Valheim",
		"script": "/path/to/Valheim/start_game_bepinex.sh"
	}
