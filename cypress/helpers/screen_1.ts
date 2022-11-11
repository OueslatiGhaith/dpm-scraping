export default class Screen1 {
  static getGourvernourats() {
    return cy.get("body").find("select[name='cod_gouv']").find("option");
  }

  static selectGouvernourat(name: string) {
    cy.get("body").find("select[name='cod_gouv']").select(name.trim());
  }

  static selectJourNuit(name: keyof typeof JOUR_NUIT) {
    cy.get("body").findAllByRole("radio").filter(`[value="${JOUR_NUIT[name]}"]`).check();
  }

  static continue() {
    cy.get("body").find("input[type='submit']").click();
  }
}

const JOUR_NUIT = {
  JOUR: "ON",
  NUIT: "OFF",
};
