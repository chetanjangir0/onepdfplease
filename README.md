# OnePDFPlease - A TUI for working with PDFs

<!-- ![merge](assets/merge.mp4) -->
<video src="assets/demo.mp4" controls width="700"></video>

A terminal-based PDF toolkit with vim keybindings providing a keyboard-driven
interface for various pdf releated tasks to ease the pain of working with pdfs

---

## Features

- Merge Pdfs
- Split Pdfs
- Encrypt Pdfs
- Decrypt Pdfs
- Convert images to pdf
- Extract embeded images from pdf
- Convert DOC/DOCX to pdf (requires [LibreOffice](https://www.libreoffice.org/get-help/install-howto/))
- Vim-style keybindings (`j`/`k` for navigation)
- Minimal UI with clear feedback and status

---

## Run without installation (Works Only on Nix/Nixos)

```bash
nix run github:chetanjangir0/onepdfplease
```

---

## Install on Nixos (If you use flakes)

Add below to your flake.nix

```nix
{
  inputs.onepdfplease.url = "github:chetanjangir0/onepdfplease";

  outputs = { self, nixpkgs, onepdfplease, ... }:
  {
    nixosConfigurations.mySystem = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        ({ pkgs, ... }: {
          environment.systemPackages = [
            onepdfplease.packages.${pkgs.system}.default
          ];
        })
      ];
    };
  };
}
```

---

## Manual Installation (Linux)

### Requirements for manual installation

- Go ≥ 1.18

```bash
git clone https://github.com/chetanjangir0/onepdfplease.git
cd onepdfplease
go build -o onepdfplease
sudo mv onepdfplease /usr/local/bin
```

Now run it from anywhere:

```bash
onepdfplease
```

---

<!-- ## Screenshots -->
<!---->
<!-- ### Main Menu   -->
<!-- ![Main Menu](screenshots/screenshot1.png) -->
<!---->
<!-- ### Paired Devices Menu   -->
<!-- ![Paired Devices](screenshots/screenshot2.png) -->
<!---->
<!-- ### Password Input Prompt   -->
<!-- ![Password Input](screenshots/screenshot3.png) -->
<!---->
<!-- --- -->

## Feedback

If you try this, please open an issue or discussion.
Even small UX feedback is welcome.

## License

MIT
