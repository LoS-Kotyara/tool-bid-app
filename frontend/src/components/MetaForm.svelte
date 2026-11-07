<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { MetaInfoState } from '../types';

  export let meta: MetaInfoState;

  const dispatch = createEventDispatcher<{
    change: void;
  }>();

  function onFieldChange() {
    dispatch('change');
  }

  function setToday(field: 'bidDate' | 'deliveryDate') {
    const now = new Date();
    const d = String(now.getDate()).padStart(2, '0');
    const m = String(now.getMonth() + 1).padStart(2, '0');
    const y = now.getFullYear();
    meta[field] = `${d}.${m}.${y}`;
    dispatch('change');
  }

  function clearDate(field: 'bidDate' | 'deliveryDate') {
    meta[field] = '';
    dispatch('change');
  }
</script>

<div class="meta-form">
  <div class="form-row">
    <div class="form-group col-small">
      <label for="bidNum">Номер №</label>
      <input id="bidNum" type="text" placeholder="52" bind:value={meta.bidNum} on:input={onFieldChange} />
    </div>

    <div class="form-group col-date">
      <label for="bidDate">Дата заявки</label>
      <div class="date-input-wrap">
        <input id="bidDate" type="text" placeholder="ДД.ММ.ГГГГ" bind:value={meta.bidDate} on:input={onFieldChange} />
        <div class="date-actions">
          <button type="button" class="btn-date" title="Сегодня" on:click={() => setToday('bidDate')}>С</button>
          <button type="button" class="btn-date" title="Очистить" on:click={() => clearDate('bidDate')}>×</button>
        </div>
      </div>
    </div>

    <div class="form-group col-dept">
      <label for="dept">Подразделение</label>
      <input id="dept" type="text" placeholder="Цех 3" bind:value={meta.dept} on:input={onFieldChange} />
    </div>

    <div class="form-group col-shop">
      <label for="shopName">Пункт назначения</label>
      <input id="shopName" type="text" placeholder="Механосборочный цех № 3" bind:value={meta.shopName} on:input={onFieldChange} />
    </div>
  </div>

  <div class="form-row split-half">
    <div class="form-group">
      <label for="initRole">Должность инициатора</label>
      <input id="initRole" type="text" placeholder="Инженер-программист" bind:value={meta.initRole} on:input={onFieldChange} />
    </div>
    <div class="form-group">
      <label for="initFio">ФИО инициатора</label>
      <input id="initFio" type="text" placeholder="Давиденко А.А." bind:value={meta.initFio} on:input={onFieldChange} />
    </div>
  </div>

  <div class="form-row split-half">
    <div class="form-group">
      <label for="curRole">Должность куратора</label>
      <input id="curRole" type="text" placeholder="Начальник цеха №3" bind:value={meta.curRole} on:input={onFieldChange} />
    </div>
    <div class="form-group">
      <label for="curFio">ФИО куратора</label>
      <input id="curFio" type="text" placeholder="Рыбкин П.В." bind:value={meta.curFio} on:input={onFieldChange} />
    </div>
  </div>

  <div class="form-row">
    <div class="form-group col-date">
      <label for="deliveryDate">Дата поставки</label>
      <div class="date-input-wrap">
        <input id="deliveryDate" type="text" placeholder="ДД.ММ.ГГГГ" bind:value={meta.deliveryDate} on:input={onFieldChange} />
        <div class="date-actions">
          <button type="button" class="btn-date" title="Сегодня" on:click={() => setToday('deliveryDate')}>С</button>
          <button type="button" class="btn-date" title="Очистить" on:click={() => clearDate('deliveryDate')}>×</button>
        </div>
      </div>
    </div>

    <div class="form-group col-flex-1">
      <label for="purpose">Цель закупки</label>
      <input id="purpose" type="text" placeholder="Для изготовления КЭМЗ БК8.035.825" bind:value={meta.purpose} on:input={onFieldChange} />
    </div>
  </div>

  <div class="form-row full-width">
    <div class="form-group">
      <label for="qualityReq">Требования к качеству (100% ширины)</label>
      <input id="qualityReq" type="text" placeholder="Без дефектов, заводской сертификат качества" bind:value={meta.qualityReq} on:input={onFieldChange} />
    </div>
  </div>
</div>

<style>
  .meta-form {
    padding: 8px 12px;
    background-color: var(--ctp-mantle);
    border-bottom: 1px solid var(--ctp-surface0);
    display: flex;
    flex-direction: column;
    gap: 6px;
    user-select: none;
  }
  .form-row {
    display: flex;
    gap: 8px;
    align-items: flex-end;
  }
  .split-half .form-group {
    flex: 1 1 50%;
  }
  .full-width .form-group {
    flex: 1 1 100%;
  }
  .col-small { width: 75px; flex-shrink: 0; }
  .col-date { width: 140px; flex-shrink: 0; }
  .col-dept { width: 120px; flex-shrink: 0; }
  .col-shop { flex: 1; }
  .col-flex-1 { flex: 1; }
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  label {
    font-size: 8.5pt;
    text-transform: uppercase;
    color: var(--ctp-subtext0);
    letter-spacing: 0.3px;
  }
  input {
    height: 26px;
    background-color: var(--ctp-surface0);
    border: 1px solid var(--ctp-surface1);
    color: var(--ctp-text);
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 10pt;
    outline: none;
    transition: border-color 0.15s;
    width: 100%;
  }
  input:focus {
    border-color: var(--ctp-blue);
  }
  .date-input-wrap {
    display: flex;
    position: relative;
    align-items: center;
  }
  .date-input-wrap input {
    padding-right: 38px;
  }
  .date-actions {
    position: absolute;
    right: 3px;
    display: flex;
    gap: 2px;
  }
  .btn-date {
    background: var(--ctp-surface1);
    border: none;
    color: var(--ctp-text);
    font-size: 8.5pt;
    width: 16px;
    height: 20px;
    border-radius: 2px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .btn-date:hover {
    background-color: var(--ctp-surface2);
    color: var(--ctp-blue);
  }
</style>
