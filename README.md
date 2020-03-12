# Extensions-Animal-Facts
A sample extension that demonstrates how to leverage the Extensions Configuration Service in tandem with an EBS. 

## What's in the Sample
The Animal Facts extension demonstrates a simple scenario of storing and retrieving extension configuration data. Specifically, the extension displays animal facts (ie. "Cats are unable to detect sweetness in anything they taste.") to viewers based on the animal type (cat or dog) that the broadcaster has selected when configuring the extension.


## Requirements
- Node.JS LTS+ with [`yarn`](https://yarnpkg.com/en/) for package management.
- Go 1.10+ with [`dep`](https://github.com/golang/dep) for package management. 

You can install `yarn` by running:
```bash
npm i -g yarn
```

You can install `dep` by running:
```bash
brew install dep
```

## First time Usage

The recommended path to using this sample is with the [Developer Rig](https://dev.twitch.tv/docs/extensions/rig/).

First, clone the repository into the folder of your choice. 

Next, do the following: 

1. Change directories to the cloned folder.
2. ONLY IF RUNNING LOCALLY (not in the developer rig): Run `yarn install` to install all prerequisite packages needed to run the frontend. 
3. ONLY IF RUNNING LOCALLY (not in the developer rig): Run `yarn cert` to generate the needed certificates; this allows the frontend and EBS servers to be run over HTTPS.
4. Change directories to the `/ebs` folder.
5. Open the `.env` file and set your extension `Client ID`, `Secret`, and the `User ID` of the extension owner (likely you).
6. Run `dep ensure` to install all prerequisite packages needed to run the EBS. 
7. Run `go build` to compile the EBS binary.

You can now start the sample from the Developer Rig. In the Project Overview, set the following fields:

1. Project Name: `Animal Facts`
2. Front-end Files Location: `path-to-clone-repository/`
3. Front-end Host Command: `yarn start`
4. Back-end Run Command: `ebs`
5. Project Folder: `path-to-clone-repository/ebs`

## File Structure

The file structure in the template is laid out with the following: 

### bin

The `/bin` folder holds the cert generation script. 

### conf 

The `/conf` folder holds the generated certs after the cert generation script runs. 

### dist

`/dist` holds the final JS files after building. 

### public

`/public` houses the static HTML files used for your code's entrypoint. 

### src

This folder houses all source code and relevant files (such as images). Each React class/component is given a folder to house all associated files (such as associated CSS).

Below this folder, the structure is much simpler.

This would be: 

`
components
-\App
--\App.js
--\App.test.js
--\App.css
-\Authentication
--\Authentication.js
...
`

### ebs

This folder houses all of the source code for the sample backend. 
