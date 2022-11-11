import DATASET_JOUR from "../../data/empty-jour.json";
import DATASET_NUIT from "../../data/empty-nuit.json";

import Screen1 from "../../helpers/screen_1";
import Screen2 from "../../helpers/screen_2";
import Screen3, { Officine } from "../../helpers/screen_3";

const OFFICINE = "http://www.dpm.tn/dpm_pharm/asppharm/listgouv.php";

type AllOfficines = {
  [Gouvernourat: string]: {
    [delegation: string]: Officine[];
  };
};
const all_officines: AllOfficines = {};

describe("officines JOUR", () => {
  it("TEST", () => {
    cy.visit(OFFICINE);

    const jourNuit = "JOUR";

    Screen1.getGourvernourats().each((gov, nbGov) => {
      const govText = gov.text().trim();
      cy.log(`~~~~~~~~~~~~~~~~~~~~~~~ ${govText} / ${jourNuit} / Officines ~~~~~~~~~~~~~~~~~~~~~~~`);

      cy.visit(OFFICINE);

      Screen1.selectGouvernourat(gov.text());
      Screen1.selectJourNuit(jourNuit);
      Screen1.continue();

      Screen2.getDelegations(jourNuit).each((del, nbDel) => {
        const delText = del.text().trim();
        cy.log(delText);

        cy.visit(OFFICINE);

        Screen1.selectGouvernourat(govText);
        Screen1.selectJourNuit(jourNuit);
        Screen1.continue();

        Screen2.selectDelegation(delText, jourNuit);
        Screen2.continue();

        Screen3.extractOfficines((officines) => {
          if (all_officines[govText] === undefined) {
            all_officines[govText] = {};
          }
          all_officines[govText][delText] = officines;
        });
      });
    });

    cy.writeFile(`./cypress/results/officines/${jourNuit}.json`, all_officines);
  });
});
