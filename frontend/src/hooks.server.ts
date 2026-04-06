import Bree from "bree";
import path from "path";

export const bree = new Bree({
    root: path.resolve(path.join("src", "jobs")),
    jobs: [
        {
            name: "generate-thumbnails",
            interval: "1m",
        },
    ],
    defaultExtension: 'ts',
});

(async () => {
    console.log("Starting Bree...");
  await bree.start();
})();