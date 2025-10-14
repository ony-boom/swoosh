Simple audio output switcher applet for Linux. _(because I do not use a DE)_

![swoosh](./assets/screenshot.png)

## Installation

```bash
go install github.com/ony-boom/swoosh@latest
```

**With nix**:
Using the package from the flake:

```nix
# In your flake.nix
{
  inputs.swoosh.url = "github:ony-boom/swoosh";
  # optional, if you want to use the same nixpkgs as your flake
  inputs.swoosh.inputs.nixpkgs.follows = "nixpkgs";
}

# somewhere in your packages if nixos
{
  environment.systemPackages = with pkgs; [
    inputs.swoosh.packages.${pkgs.system}.swoosh
  ];
}

# or in your home-manager configuration
{
  home.packages = with pkgs; [
    inputs.swoosh.packages.${pkgs.system}.swoosh
  ];
```

Or you can install directly like this:

```bash
# if newer version of nix
nix profile add github:ony-boom/swoosh
# or with older version
nix profile install github:ony-boom/swoosh
```

**Roadmap:**

- [ ] List available audio inputs (sources)
- [ ] Allow simple configuration (e.g. polling interval, hide source, ...)
- [ ] Better way to detect signal changes (currently polling) or using a different tray library that allow rerendering every time the menu is opened

**Alternatives:**

- [pasystray](https://github.com/christophgysin/pasystray)
- [indicator-sound-switcher](https://github.com/yktoo/indicator-sound-switcher)
