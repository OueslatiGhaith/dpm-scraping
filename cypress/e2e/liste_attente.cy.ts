import Screen1 from "../helpers/screen_1";
import Screen2 from "../helpers/screen_2";
import Screen3, { Attente } from "../helpers/screen_3";

const LISTE_ATTENTE = "http://www.dpm.tn/dpm_pharm/asppharm/listgouv_att.php";
const OFFICINE = "http://www.dpm.tn/dpm_pharm/asppharm/listgouv.php";

type ListAttente = {
  [Gouvernourat: string]: {
    [delegation: string]: Attente;
  };
};
const list_attente: ListAttente = {};

// describe("TEST", () => {
//   it("passes", () => {
//     cy.visit(LISTE_ATTENTE);

//     const jourNuit = "JOUR";
//     const gouvernourat = "BEJA";
//     const delegation = "GOUBELLAT";

//     Screen1.selectGouvernourat(gouvernourat);
//     Screen1.selectJourNuit(jourNuit);
//     Screen1.continue();

//     Screen2.selectDelegation(delegation, jourNuit);
//     Screen2.continue();

//     Screen3.extractAttente((attente) => {
//       if (list_attente[gouvernourat] === undefined) {
//         list_attente[gouvernourat] = {};
//       }
//       list_attente[gouvernourat][delegation] = attente;
//     });

//     cy.writeFile(`./cypress/results/liste_attente/${jourNuit}.json`, list_attente);
//   });
// });

describe("liste attente JOUR", () => {
  it("passes", () => {
    cy.visit(LISTE_ATTENTE);

    const jourNuit = "JOUR";

    Screen1.getGourvernourats().each((gov, nbGov) => {
      const govText = gov.text().trim();
      cy.log(`~~~~~~~~~~~~~~~~~~~~~~~ ${govText} / ${jourNuit} / Liste d'attente ~~~~~~~~~~~~~~~~~~~~~~~`);
      cy.visit(LISTE_ATTENTE);

      Screen1.selectGouvernourat(gov.text());
      Screen1.selectJourNuit(jourNuit);
      Screen1.continue();

      Screen2.getDelegations(jourNuit).each((del, nbDel) => {
        const delText = del.text().trim();
        cy.log(delText);
        cy.visit(LISTE_ATTENTE);

        Screen1.selectGouvernourat(govText);
        Screen1.selectJourNuit(jourNuit);
        Screen1.continue();

        Screen2.selectDelegation(delText, jourNuit);
        Screen2.continue();

        Screen3.extractAttente((attente) => {
          if (list_attente[govText] === undefined) {
            list_attente[govText] = {};
          }
          list_attente[govText][delText] = attente;
        });
      });
    });

    cy.writeFile(`./cypress/results/liste_attente/${jourNuit}.json`, list_attente);
  });
});

describe("liste attente NUIT", () => {
  it("passes", () => {
    cy.visit(LISTE_ATTENTE);

    const jourNuit = "NUIT";

    Screen1.getGourvernourats().each((gov, nbGov) => {
      const govText = gov.text().trim();
      cy.log(`~~~~~~~~~~~~~~~~~~~~~~~ ${govText} / ${jourNuit} / Liste d'attente ~~~~~~~~~~~~~~~~~~~~~~~`);
      cy.visit(LISTE_ATTENTE);

      Screen1.selectGouvernourat(gov.text());
      Screen1.selectJourNuit(jourNuit);
      Screen1.continue();

      Screen2.getDelegations(jourNuit).each((del, nbDel) => {
        const delText = del.text().trim();
        cy.log(delText);
        cy.visit(LISTE_ATTENTE);

        Screen1.selectGouvernourat(govText);
        Screen1.selectJourNuit(jourNuit);
        Screen1.continue();

        Screen2.selectDelegation(delText, jourNuit);
        Screen2.continue();

        Screen3.extractAttente((attente) => {
          if (list_attente[govText] === undefined) {
            list_attente[govText] = {};
          }
          list_attente[govText][delText] = attente;
        });
      });
    });

    cy.writeFile(`./cypress/results/liste_attente/${jourNuit}.json`, list_attente);
  });
});
