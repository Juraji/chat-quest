import {ResolveFn} from '@angular/router';
import {NEW_SCENARIO, Scenario, Scenarios} from '@db/scenarios';
import {NewRecord} from '@db/core';
import {inject} from '@angular/core';


export const editScenarioResolver: ResolveFn<Scenario | NewRecord<Scenario>> = (route) => {
  const service = inject(Scenarios)
  const scenarioId = route.paramMap.get('scenarioId')!!
  const iScenarioId = Number(scenarioId)

  if (scenarioId === 'new') {
    const sceneDescription: string = route.queryParamMap.has('sceneDescription')
      ? route.queryParamMap.get('sceneDescription')!
      : NEW_SCENARIO.sceneDescription

    return {...NEW_SCENARIO, sceneDescription}
  } else if (!isNaN(iScenarioId)) {
    return service.get(iScenarioId)
  } else {
    throw new Error(`Scenario with id "${scenarioId}" not found`)
  }
}
