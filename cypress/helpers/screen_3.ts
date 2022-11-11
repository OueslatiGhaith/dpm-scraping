export default class Screen3 {
  static extractNbOfficines(callback: (nb: string) => void) {
    cy.get("body")
      .find("font[color='#000000']")
      .last()
      .then((el) => {
        callback(el.text());
      });
  }

  static extractOfficines(callback: (officines: Officine[]) => void) {
    const data: Officine[] = [];

    cy.get("body").then(($body) => {
      const $font = $body.find("font");
      if ($font && $font.text().trim() === "Pas d'inscription sur cette liste") {
        callback(data);
      } else {
        cy.get("body")
          .find("table")
          .last()
          .find("tr")
          .each((tr, i) => {
            if (i === 0) return;

            tr.find("td").each((j, td) => {
              if (j === 0) {
                data.push({
                  ordre: i,
                  nom: td.children[0].textContent,
                  adresse: "",
                  telephone: "",
                });
              }
              if (j === 1) {
                data[data.length - 1].adresse = td.children[0].textContent;
              }
              if (j === 2) {
                data[data.length - 1].telephone = td.children[0].textContent;
              }
            });
          });

        callback(data);
      }
    });
  }

  static extractAttente(callback: (attente: Attente) => void) {
    const data: Attente = {
      zone: "",
      population: "",
      nb_officines: "",
      liste: [],
    };

    cy.get("body").then(($body) => {
      const $font = $body.find("font");
      if ($font && $font.text().trim() === "Pas d'inscription sur cette liste") {
        callback(data);
      } else {
        cy.xpath("html/body/table/tbody/tr[1]/td/p/font/b/b/font[2]").then((el) => {
          data.zone = el.text();
        });
        cy.xpath("/html/body/table/tbody/tr[1]/td/p/font/b/b/b/font[2]").then((el) => {
          data.population = el.text();
        });
        cy.xpath("/html/body/table/tbody/tr[1]/td/p/font/b/b/b/b/font[2]").then((el) => {
          data.nb_officines = el.text();
        });

        cy.get("body")
          .find("table")
          .last()
          .find("tr")
          .each((tr, i) => {
            if (i === 0) return;

            tr.find("td").each((j, td) => {
              if (j === 0) {
                data.liste.push({
                  ordre: td.children[0].textContent,
                  nom: "",
                  prenom: "",
                  epouse: "",
                  date_inscription: "",
                });
              }
              if (j === 1) {
                data.liste[data.liste.length - 1].nom = td.children[0].textContent;
              }
              if (j === 2) {
                data.liste[data.liste.length - 1].prenom = td.children[0].textContent;
              }
              if (j === 3) {
                data.liste[data.liste.length - 1].epouse = td.children[0].textContent;
              }
              if (j === 4) {
                data.liste[data.liste.length - 1].date_inscription = td.children[0].textContent;
              }
            });
          });

        callback(data);
      }
    });

    console.log(data);
  }
}

export type Officine = {
  ordre: number;
  nom: string;
  adresse: string;
  telephone: string;
};

export type Attente = {
  zone: string;
  population: string;
  nb_officines: string;
  liste: {
    ordre: string;
    nom: string;
    prenom: string;
    epouse: string;
    date_inscription: string;
  }[];
};
