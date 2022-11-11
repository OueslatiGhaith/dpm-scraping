export default class Screen2 {
  static getDelegations(time: "JOUR" | "NUIT") {
    const input_name = time === "JOUR" ? "cod_del" : "cod_comm";
    const delegations: string[] = [];
    return cy.get("body").find(`select[name='${input_name}']`).find("option");
  }

  static selectDelegation(name: string, time: "JOUR" | "NUIT") {
    const input_name = time === "JOUR" ? "cod_del" : "cod_comm";
    cy.get("body").find(`select[name='${input_name}']`).select(name);
  }

  static continue() {
    cy.get("body").find("input[type='submit']").click();
  }
}
