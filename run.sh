#! /bin/bash
pnpm cypress:run --spec cypress/e2e/ghazi.cy.ts
node ./extract.js > a.csv