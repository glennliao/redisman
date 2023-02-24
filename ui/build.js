const fs = require("fs");
const fsExtra = require("fs-extra");

const desktopDist = "../cmd/desktop/dist";
fs.rmSync(desktopDist, { force: true, recursive: true });
fsExtra.copy("./dist", desktopDist);

