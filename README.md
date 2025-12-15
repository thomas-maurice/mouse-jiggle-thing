# Mouse Jiggle Thing 🐭

## English

### Is your boss a cuntish idiot who wants to track your every movement down to your mouse clicks?

**Good news!** This project makes that bullshit stop working.

The preferred kind of device to achieve such shenanigans is [this one](https://nl.aliexpress.com/w/wholesale-TENSTAR-RP2350%2525252dUSB.html)

### What does this do?

Your company installed some dystopian surveillance software that tracks mouse movements to see if you're "actually working"? This firmware turns a Raspberry Pi Pico into a fake mouse that jiggles randomly, making their tracking systems think you're perpetually busy.

The best part? **It's completely undetectable.** Your computer just sees a legitimate mouse device. No suspicious software to install, no traces to find. Just a tiny board plugged into USB that keeps your workstation "active" while you go grab a coffee, take a shit, or do literally anything more valuable than being monitored like a prisoner.

### How it works

1. Flash a firmware onto your Raspberry Pi Pico (or Pico 2). You can choose from like 803 different models so it is likely that your corporate provided one is included.
2. The Pico pretends to be any mouse from our database (803 different models!)
3. Plug it in, and it jiggles the mouse cursor periodically (with some randomness)
4. Your paranoid tracking software thinks you're working
5. You're free to live your life

### Features

- **803 different mouse firmwares** - Pick any mouse model you want
- **Works on RP2040 and RP2350** - Supports both Raspberry Pi Pico and Pico 2
- **Completely passive** - No software to install on your computer
- **Undetectable** - Shows up as a real USB mouse device
- **Open source** - Audit the code, modify it, do whatever you want

### Usage

1. Download `firmwares.zip` from the [latest release](https://github.com/thomas-maurice/mouse-jiggle-thing/releases)
2. Extract it and choose your board folder (`rp2040/` or `rp2350/`)
3. Pick a mouse firmware (they're named by vendor, product, and VID.PID)
4. Hold the BOOTSEL button on your Pico while plugging it in
5. Copy the `.uf2` file to the `RPI-RP2` drive that appears
6. Done! Your Pico is now a mouse jiggler

### Building locally

```bash
# Build all firmwares for both RP2040 and RP2350
./scripts/generate_all_firmwares.sh

# Build only for RP2040
TARGET=rp2040 ./scripts/generate_all_firmwares.sh

# Build only for RP2350
TARGET=rp2350 ./scripts/generate_all_firmwares.sh
```

---

## Français

### Votre patron est un connard qui veut traquer chacun de vos mouvements jusqu'aux clics de souris ?

**Bonne nouvelle !** Ce projet fait en sorte que ces conneries ne marchent plus.

### Appareils supportés
Les trucs comme [ça](https://nl.aliexpress.com/w/wholesale-TENSTAR-RP2350%2525252dUSB.html)

### Qu'est-ce que ça fait ?

Votre boîte a installé un logiciel de surveillance dystopique qui traque les mouvements de souris pour voir si vous "travaillez vraiment" ? Ce firmware transforme votre Raspberry Pi Pico en fausse souris qui bouge aléatoirement, faisant croire à leur système de tracking que vous êtes perpétuellement occupé.

Le meilleur ? **C'est complètement indétectable.** Votre ordinateur voit juste un périphérique souris légitime. Pas de logiciel suspect à installer, aucune trace à trouver. Juste une petite carte branchée en USB qui garde votre poste "actif" pendant que vous allez chercher un café, prendre une chiasse, ou faire littéralement n'importe quoi de plus utile que d'être surveillé comme un prisonnier.

### Comment ça marche

1. Flashez un firmware sur votre Raspberry Pi Pico (ou Pico 2)
2. Le Pico se fait passer pour n'importe quelle souris de notre base de données (803 modèles différents !)
3. Branchez-le, et il fait bouger le curseur de souris périodiquement
4. Votre logiciel de tracking parano pense que vous travaillez
5. Vous êtes libre de vivre votre vie

### Fonctionnalités

- **803 firmwares de souris différents** - Choisissez le modèle de souris que vous voulez
- **Fonctionne sur RP2040 et RP2350** - Support du Raspberry Pi Pico et Pico 2
- **Complètement passif** - Aucun logiciel à installer sur votre ordinateur
- **Indétectable** - Apparaît comme un vrai périphérique USB souris
- **Open source** - Auditez le code, modifiez-le, faites ce que vous voulez

### Utilisation

1. Téléchargez `firmwares.zip` depuis la [dernière release](https://github.com/thomas-maurice/mouse-jiggle-thing/releases)
2. Extrayez-le et choisissez le dossier de votre carte (`rp2040/` ou `rp2350/`)
3. Choisissez un firmware de souris (ils sont nommés par fabricant, produit, et VID.PID)
4. Maintenez le bouton BOOTSEL de votre Pico en le branchant
5. Copiez le fichier `.uf2` sur le lecteur `RPI-RP2` qui apparaît
6. C'est fait ! Votre Pico est maintenant un jiggleur de souris

### Compiler localement

```bash
# Compiler tous les firmwares pour RP2040 et RP2350
./scripts/generate_all_firmwares.sh

# Compiler seulement pour RP2040
TARGET=rp2040 ./scripts/generate_all_firmwares.sh

# Compiler seulement pour RP2350
TARGET=rp2350 ./scripts/generate_all_firmwares.sh
```

---

## License

Do whatever the fuck you want with this. Seriously.

## Contributing

PRs welcome. Let's make workplace surveillance harder, one jiggle at a time.

## Disclaimer

This project is for educational purposes and to protest against invasive workplace monitoring. Use responsibly. Don't be a dick. If your employer is tracking your mouse movements, maybe it's time to find a better job where they trust their employees.
