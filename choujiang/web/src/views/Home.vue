
&lt;template&gt;
  &lt;div class="home"&gt;
    &lt;div class="header"&gt;
      &lt;h1&gt;🎉 幸运大转盘 🎉&lt;/h1&gt;
      &lt;p&gt;好运转转转，大奖抽不停！&lt;/p&gt;
    &lt;/div&gt;

    &lt;div class="wheel-container"&gt;
      &lt;div class="wheel-wrapper"&gt;
        &lt;div class="pointer"&gt;&lt;/div&gt;
        &lt;canvas ref="wheelCanvas" :style="{ transform: `rotate(${rotation}deg)` }" class="wheel"&gt;&lt;/canvas&gt;
        &lt;div class="center-btn" @click="startLottery" :class="{ disabled: isSpinning }"&gt;
          {{ isSpinning ? '抽奖中...' : '开始抽奖' }}
        &lt;/div&gt;
      &lt;/div&gt;
    &lt;/div&gt;

    &lt;div class="footer"&gt;
      &lt;el-button type="primary" class="my-prizes-btn" @click="$router.push('/my-prizes')"&gt;
        🎁 我的奖品
      &lt;/el-button&gt;
    &lt;/div&gt;

    &lt;el-dialog v-model="userFormVisible" title="请填写信息" width="90%"&gt;
      &lt;el-form :model="userForm" label-width="80px"&gt;
        &lt;el-form-item label="姓名"&gt;
          &lt;el-input v-model="userForm.name" placeholder="请输入姓名"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
        &lt;el-form-item label="手机号"&gt;
          &lt;el-input v-model="userForm.phone" placeholder="请输入11位手机号" maxlength="11"&gt;&lt;/el-input&gt;
        &lt;/el-form-item&gt;
      &lt;/el-form&gt;
      &lt;template #footer&gt;
        &lt;el-button @click="userFormVisible = false"&gt;取消&lt;/el-button&gt;
        &lt;el-button type="primary" @click="submitUserForm"&gt;确定&lt;/el-button&gt;
      &lt;/template&gt;
    &lt;/el-dialog&gt;

    &lt;el-dialog v-model="resultVisible" :title="resultTitle" width="90%"&gt;
      &lt;div class="result-content"&gt;
        &lt;div class="result-icon"&gt;{{ resultIcon }}&lt;/div&gt;
        &lt;div class="result-text"&gt;{{ resultText }}&lt;/div&gt;
      &lt;/div&gt;
      &lt;template #footer&gt;
        &lt;el-button type="primary" @click="resultVisible = false"&gt;确定&lt;/el-button&gt;
      &lt;/template&gt;
    &lt;/el-dialog&gt;
  &lt;/div&gt;
&lt;/template&gt;

&lt;script setup&gt;
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getPrizes, doLottery } from '../utils/api'

const router = useRouter()
const wheelCanvas = ref(null)
const rotation = ref(0)
const isSpinning = ref(false)
const prizes = ref([])
const userFormVisible = ref(false)
const resultVisible = ref(false)
const resultTitle = ref('')
const resultIcon = ref('')
const resultText = ref('')

const userForm = ref({
  name: '',
  phone: ''
})

const colors = [
  '#FF6B6B',
  '#4ECDC4',
  '#FFE66D',
  '#95E1D3',
  '#F38181',
  '#AA96DA',
  '#FCBAD3',
  '#A8D8EA'
]

onMounted(() =&gt; {
  loadPrizes()
  const saved = localStorage.getItem('lottery_user')
  if (saved) {
    userForm.value = JSON.parse(saved)
  }
})

const loadPrizes = async () =&gt; {
  try {
    const res = await getPrizes()
    prizes.value = res.data.slice(0, 8)
    if (prizes.value.length &lt; 8) {
      const addCount = 8 - prizes.value.length
      for (let i = 0; i &lt; addCount; i++) {
        prizes.value.push({ id: 0, name: '谢谢参与', probability: 0, enabled: true })
      }
    }
    nextTick(() =&gt; drawWheel())
  } catch (err) {
    ElMessage.error('加载奖品失败')
  }
}

const drawWheel = () =&gt; {
  const canvas = wheelCanvas.value
  canvas.width = 320
  canvas.height = 320
  const ctx = canvas.getContext('2d')
  const centerX = canvas.width / 2
  const centerY = canvas.height / 2
  const radius = 150

  const segmentAngle = (2 * Math.PI) / prizes.value.length

  prizes.value.forEach((prize, index) =&gt; {
    ctx.beginPath()
    ctx.moveTo(centerX, centerY)
    ctx.arc(centerX, centerY, radius, index * segmentAngle, (index + 1) * segmentAngle)
    ctx.closePath()
    ctx.fillStyle = colors[index % colors.length]
    ctx.fill()
    ctx.strokeStyle = '#fff'
    ctx.lineWidth = 3
    ctx.stroke()

    ctx.save()
    ctx.translate(centerX, centerY)
    ctx.rotate(index * segmentAngle + segmentAngle / 2)
    ctx.textAlign = 'center'
    ctx.fillStyle = '#333'
    ctx.font = 'bold 14px Arial'
    ctx.fillText(prize.name, radius / 2, 5)
    ctx.restore()
  })
}

const startLottery = () =&gt; {
  if (isSpinning.value) return

  if (!userForm.value.name || !userForm.value.phone) {
    userFormVisible.value = true
    return
  }

  if (!/^1\d{10}$/.test(userForm.value.phone)) {
    ElMessage.warning('请输入正确的11位手机号')
    return
  }

  isSpinning.value = true
  doLottery({ name: userForm.value.name, phone: userForm.value.phone }).then(res =&gt; {
    const data = res.data
    const prizeIndex = prizes.value.findIndex(p =&gt; p.name === data.prizeName)
    const targetIndex = prizeIndex &gt;= 0 ? prizeIndex : 0

    const spins = 5
    const segmentAngle = 360 / prizes.value.length
    const targetAngle = 360 - (targetIndex * segmentAngle + segmentAngle / 2)
    const totalRotation = rotation.value + spins * 360 + targetAngle - (rotation.value % 360)

    rotation.value = totalRotation

    setTimeout(() =&gt; {
      isSpinning.value = false
      if (data.isWin) {
        resultTitle.value = '🎉 恭喜中奖！'
        resultIcon.value = '🎁'
        resultText.value = `恭喜您获得：${data.prizeName}`
      } else {
        resultTitle.value = '😊 谢谢参与'
        resultIcon.value = '🍀'
        resultText.value = '很遗憾，下次一定会中奖的！'
      }
      resultVisible.value = true
    }, 4000)
  }).catch(() =&gt; {
    isSpinning.value = false
    ElMessage.error('抽奖失败，请稍后重试')
  })
}

const submitUserForm = () =&gt; {
  if (!userForm.value.name) {
    ElMessage.warning('请输入姓名')
    return
  }
  if (!/^1\d{10}$/.test(userForm.value.phone)) {
    ElMessage.warning('请输入正确的11位手机号')
    return
  }
  localStorage.setItem('lottery_user', JSON.stringify(userForm.value))
  userFormVisible.value = false
  startLottery()
}
&lt;/script&gt;

&lt;style scoped&gt;
.home {
  min-height: 100vh;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.header {
  text-align: center;
  margin-bottom: 30px;
  color: #fff;
}

.header h1 {
  font-size: 28px;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.2);
}

.header p {
  font-size: 16px;
  opacity: 0.9;
}

.wheel-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wheel-wrapper {
  position: relative;
  width: 320px;
  height: 320px;
}

.wheel {
  width: 320px;
  height: 320px;
  border-radius: 50%;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  transition: transform 4s cubic-bezier(0.2, 0.8, 0.3, 1);
}

.pointer {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 20px solid transparent;
  border-right: 20px solid transparent;
  border-top: 40px solid #ff4757;
  z-index: 10;
  filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.3));
}

.center-btn {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ff6b6b 0%, #ff4757 100%);
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 6px 20px rgba(255, 71, 87, 0.5);
  z-index: 20;
  user-select: none;
  transition: transform 0.2s, box-shadow 0.2s;
}

.center-btn:hover:not(.disabled) {
  transform: translate(-50%, -50%) scale(1.05);
  box-shadow: 0 8px 25px rgba(255, 71, 87, 0.6);
}

.center-btn:active:not(.disabled) {
  transform: translate(-50%, -50%) scale(0.95);
}

.center-btn.disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.footer {
  padding: 20px 0;
}

.my-prizes-btn {
  font-size: 16px;
  padding: 12px 30px;
  border-radius: 25px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.result-content {
  text-align: center;
  padding: 20px 0;
}

.result-icon {
  font-size: 60px;
  margin-bottom: 20px;
}

.result-text {
  font-size: 20px;
  color: #333;
}
&lt;/style&gt;
