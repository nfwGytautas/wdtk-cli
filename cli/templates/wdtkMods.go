package templates

// PUBLIC TYPES
// ========================================================================

// WDTK services build template
const WDTKBuildTemplate = `
cd ..

if [ ! -d 'wdtk-services' ]
then
    echo Cloning modules
    git clone https://github.com/nfwGytautas/wdtk-services
fi

cd wdtk-services

echo Updating
git pull

. BUILD.sh

# Return to deploy
cd ..

echo $(pwd)
echo Copying
cp wdtk-services/bin/* bin/unix/

`
