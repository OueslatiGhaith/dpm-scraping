describe("empty spec", () => {
  it("passes", () => {
    cy.visit("http://www.sciencespharmaceutiques.org.tn/fr/revue_resultat.php?action=6");
    cy.get("a").contains("Revue Essaydali").click();
    const data: string[] = [];
    cy.get("select")
      .first()
      .children()
      .each((option) => {
        data.push(option.text());
      });

    const oldData = cy.readFile("./cypress/results/ghazi.json");
    oldData.then((content: string[]) => {
      const isSame = content.every((item) => data.includes(item));

      if (!isSame) {
        cy.writeFile(
          "/home/go/Desktop/GOOGLE HAS SPIED ON YOU, PLEASE SEND 10 BTC.txt",
          `
          pwease visit http://www.sciencespharmaceutiques.org.tn/fr/revue_resultat.php?action=6
          sank you!
        `
        );
      }
    });
  });
});
