sed -n '/id="icon--/p' ./icons/sprite/sprite.vue > icons/icon-names.ts
sed -i'' -e 's/.*<symbol id="icon--/  | "/g' icons/icon-names.ts
sed -i'' -e 's/" viewBox=.*/"/g' icons/icon-names.ts
sed -i'' -e '1s/^/export type IconNames =/' icons/icon-names.ts
rm -rf icons/icon-names.ts-e
echo 'Icon names collected!'
