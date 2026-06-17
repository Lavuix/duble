const fs = require("fs");

const iconTypesPath = "./icons/icon-names.ts";
const packagePath = "../../../package.json";

const log = text => {
  const logFilePath = "./logs.txt";
  if (fs.existsSync(logFilePath)) {
    fs.appendFileSync(logFilePath, text + "\n");
  } else {
    fs.writeFileSync(logFilePath, text + "\n", { flag: "w" });
  }
};

log("run install script");

/**
 *
 * @param { string[] } typesArray
 */
const addTypes = typesArray => {
  fs.readFile(iconTypesPath, "utf8", (err, data) => {
    let splitArray = data.split("\n");
    const content =
      "  | " + typesArray.map(n => `"${n}"`).join("\n  | ") + "\n;";

    const semicolonIndex = splitArray.findIndex(line => line.includes(";"));

    if (semicolonIndex !== -1) {
      splitArray.splice(
        semicolonIndex,
        splitArray.length - semicolonIndex,
        content
      );
    } else {
      splitArray.push(content);
    }

    fs.writeFileSync(iconTypesPath, splitArray.join("\n") + "\n");
  });

  const content = "  | " + typesArray.map(n => `"${n}"`).join("\n  | ") + "\n;";
  fs.appendFile(iconTypesPath, content, () => {});
};

if (fs.existsSync(packagePath)) {
  log("package.json exists");

  /**
   * @type { { "tn-ui-kit"?: { "external-icon-names"?: string[] } } }
   */
  const parentPackage = require(packagePath);

  log("load package.json: \n\n" + JSON.stringify(parentPackage));

  /**
   *
   * @type {{"external-icon-names"?: string[]}}
   */
  const packageSettings = parentPackage["tn-ui-kit"];

  if (packageSettings) {
    log("found tn-ui-kit settings: \n\n" + JSON.stringify(packageSettings));
    const externalIconNames = packageSettings["external-icon-names"];
    if (externalIconNames && externalIconNames.length) {
      log("got externalIconNames: \n\n" + JSON.stringify(externalIconNames));
      addTypes(externalIconNames);
    }
  }
} else {
  log("could not found package.json");
}
