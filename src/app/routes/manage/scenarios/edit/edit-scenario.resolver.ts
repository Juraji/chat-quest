import {ResolveFn} from '@angular/router';
import {Scenario, Scenarios} from '@db/scenarios';
import {NewRecord} from '@db/core';
import {inject} from '@angular/core';

const NEW_SCENARIO: NewRecord<Scenario> = {
  name: '',
  sceneDescription: ''
}

export const editScenarioResolver: ResolveFn<Scenario | NewRecord<Scenario>> = (route) => {
  const service = inject(Scenarios)
  const scenarioId = route.paramMap.get('scenarioId')!!
  const iScenarioId = Number(scenarioId)

  if (scenarioId === 'new') {
    return {...NEW_SCENARIO}
  } else if (!isNaN(iScenarioId)) {
    return service.get(iScenarioId)
  } else {
    throw new Error(`Scenario with id "${scenarioId}" not found`)
  }
}
