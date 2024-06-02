const fs = require("fs");

const json = fs.readFileSync(
  "./cypress/results/liste_attente/JOUR.json",
  "utf8"
);
const DATA = JSON.parse(json);

const csv = ["Governorat,Delegation,nb Officines,nb Liste"];

Object.entries(DATA).forEach(([gov, v_gov]) => {
  Object.entries(v_gov).forEach(([del, v_del]) => {
    const nb_officines = v_del.nb_officines;
    const nb_liste = v_del.liste.length;

    csv.push(`${gov},${del},${nb_officines},${nb_liste}`);
  });
});

const result = csv.join("\n");

console.log(result);
