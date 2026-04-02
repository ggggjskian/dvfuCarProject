<template>
  <div class="create-view">
    <h2>🚗 Создать поездку</h2>
    <form @submit.prevent="createTrip">
      <div class="form-group">
        <label>Тип поездки</label>
        <div class="trip-type">
          <label>
            <input type="radio" value="from_campus" v-model="form.trip_type" />
            🏫 Из кампуса
          </label>
          <label>
            <input type="radio" value="to_campus" v-model="form.trip_type" />
            📍 В кампус
          </label>
        </div>
      </div>

      <div class="form-group">
        <label>Куда / Откуда (выберите на карте)</label>
        <LocationPicker @select="onLocationSelect" />
        <div v-if="form.point" class="selected-point">
          ✅ Выбрано: {{ form.point }}
        </div>
      </div>

      <div class="form-row">
        <div class="form-group half">
          <label>Дата и время</label>
          <input
            type="datetime-local"
            id="departure_time"
            name="departure_time"
            v-model="form.departure_time"
            required
          />
        </div>
        <div class="form-group half">
          <label>Цена (₽)</label>
          <input
            type="number"
            id="price"
            name="price"
            v-model="form.price"
            min="0"
            step="50"
            required
          />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group half">
          <label>Мест всего</label>
          <input
            type="number"
            id="seats_total"
            name="seats_total"
            v-model="form.seats_total"
            min="1"
            max="8"
            required
          />
        </div>
        <div class="form-group half">
          <label>Готов отклониться (км)</label>
          <input
            type="number"
            id="max_deviation_km"
            name="max_deviation_km"
            v-model="form.max_deviation_km"
            min="0"
            max="20"
            value="3"
          />
        </div>
      </div>

      <div class="form-group">
        <label>Комментарий (необязательно)</label>
        <textarea
          id="comment"
          name="comment"
          v-model="form.comment"
          rows="2"
        ></textarea>
      </div>

      <div class="form-group">
        <label>Гибкость по времени (± минут)</label>
        <input
          type="number"
          id="time_flexibility"
          name="time_flexibility"
          v-model="form.time_flexibility_minutes"
          min="0"
          max="120"
          value="30"
        />
      </div>

      <button type="submit" class="submit-btn" :disabled="isSubmitting">
        {{ isSubmitting ? "Создание..." : "Создать поездку" }}
      </button>
    </form>

    <div v-if="error" class="error-message">{{ error }}</div>
    <div v-if="success" class="success-message">
      ✅ Поездка создана!
      <router-link :to="`/trip/${createdTripId}`">Посмотреть</router-link>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
import { useRouter } from "vue-router";
import LocationPicker from "../components/LocationPicker.vue";
import { tripsAPI } from "../services/api";

export default {
  name: "CreateView",
  components: { LocationPicker },
  setup() {
    const router = useRouter();
    const isSubmitting = ref(false);
    const error = ref("");
    const success = ref(false);
    const createdTripId = ref(null);

    const form = ref({
      trip_type: "from_campus",
      point: "",
      point_lat: null,
      point_lon: null,
      departure_time: "",
      seats_total: 4,
      price: 200,
      comment: "",
      max_deviation_km: 3,
      time_flexibility_minutes: 30,
    });

    const onLocationSelect = (location) => {
      form.value.point = location.address;
      form.value.point_lat = location.lat;
      form.value.point_lon = location.lon;
    };

    const createTrip = async () => {
      // Проверка обязательных полей
      if (!form.value.point_lat || !form.value.point_lon) {
        error.value = "Выберите точку на карте";
        return;
      }
      if (!form.value.departure_time) {
        error.value = "Укажите дату и время";
        return;
      }

      isSubmitting.value = true;
      error.value = "";

      // Формируем данные для отправки (приводим числа)
      const tripData = {
        trip_type: form.value.trip_type,
        point: form.value.point,
        point_lat: form.value.point_lat,
        point_lon: form.value.point_lon,
        departure_time: new Date(form.value.departure_time).toISOString(),
        seats_total: Number(form.value.seats_total),
        price: Number(form.value.price),
        comment: form.value.comment || null,
        max_deviation_km: Number(form.value.max_deviation_km),
        time_flexibility_minutes: Number(form.value.time_flexibility_minutes),
      };

      console.log("Sending trip data:", JSON.stringify(tripData));

      try {
        const response = await tripsAPI.create(tripData);
        console.log('Response status:', response.status);
        console.log('Response data:', response.data);
        createdTripId.value = response.data.id;
        success.value = true;
        setTimeout(() => {
          router.push(`/trip/${createdTripId.value}`);
        }, 2000);
      } catch (err) {
        console.error("Creation error:", err);
        if (err.response) {
          error.value = JSON.stringify(err.response.data.detail) || "Ошибка сервера";
        } else {
          error.value = err.message || "Ошибка создания поездки";
        }
      } finally {
        isSubmitting.value = false;
      }
    };

    return {
      form,
      isSubmitting,
      error,
      success,
      createdTripId,
      onLocationSelect,
      createTrip,
    };
  },
};
</script>

<style scoped>
.create-view {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
  padding-bottom: 80px;
}

h2 {
  margin-bottom: 24px;
  color: #333;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #555;
}

.form-row {
  display: flex;
  gap: 16px;
}

.half {
  flex: 1;
}

input[type="text"],
input[type="number"],
input[type="datetime-local"],
textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
}

.trip-type {
  display: flex;
  gap: 20px;
  padding: 8px 0;
}

.trip-type label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
}

.selected-point {
  margin-top: 8px;
  padding: 12px;
  background: #e8f4fd;
  border-radius: 8px;
  color: #0056b3;
}

.submit-btn {
  width: 100%;
  padding: 16px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: #218838;
}

.submit-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.error-message {
  margin-top: 16px;
  padding: 12px;
  background: #f8d7da;
  color: #721c24;
  border-radius: 8px;
}

.success-message {
  margin-top: 16px;
  padding: 12px;
  background: #d4edda;
  color: #155724;
  border-radius: 8px;
}
</style>