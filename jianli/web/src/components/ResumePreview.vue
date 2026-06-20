<template>
  <div class="resume-preview">
    <div class="preview-header">
      <span class="title">实时预览</span>
    </div>
    <div class="preview-body">
      <div class="a4-container" :style="containerStyle">
        <div class="resume-content" :style="contentStyle" id="resume-content">
          <template v-for="module in visibleModules" :key="module.id">
            <div v-if="module.type === 'basic'" class="module basic-module">
              <h1 class="name" :style="{ color: primaryColor }">{{ basicInfo.name || '姓名' }}</h1>
              <div class="basic-info-row">
                <span v-if="basicInfo.phone" class="info-item">📞 {{ basicInfo.phone }}</span>
                <span v-if="basicInfo.email" class="info-item">📧 {{ basicInfo.email }}</span>
                <span v-if="basicInfo.city" class="info-item">📍 {{ basicInfo.city }}</span>
                <span v-if="basicInfo.job_intention" class="info-item">💼 {{ basicInfo.job_intention }}</span>
              </div>
              <div class="basic-info-row">
                <a v-if="basicInfo.github" :href="basicInfo.github" target="_blank" class="link info-item">
                  🔗 GitHub
                </a>
                <a v-if="basicInfo.website" :href="basicInfo.website" target="_blank" class="link info-item">
                  🌐 {{ basicInfo.website }}
                </a>
              </div>
            </div>

            <div v-else-if="module.type === 'education'" class="module">
              <h2 class="module-title" :style="{ borderColor: primaryColor }">
                <span :style="{ color: primaryColor }">{{ module.name }}</span>
              </h2>
              <div v-for="(item, index) in getModuleContent(module.type)" :key="index" class="item">
                <div class="item-header">
                  <span class="item-title">{{ item.school || '学校名称' }}</span>
                  <span class="item-date">
                    {{ item.start_date || '开始' }} - {{ item.end_date || '结束' }}
                  </span>
                </div>
                <div class="item-subtitle">
                  {{ item.degree || '学历' }} · {{ item.major || '专业' }}
                </div>
                <p v-if="item.description" class="item-description">{{ item.description }}</p>
              </div>
              <el-empty v-if="getModuleContent(module.type).length === 0" description="暂无教育经历" :image-size="60" />
            </div>

            <div v-else-if="module.type === 'experience'" class="module">
              <h2 class="module-title" :style="{ borderColor: primaryColor }">
                <span :style="{ color: primaryColor }">{{ module.name }}</span>
              </h2>
              <div v-for="(item, index) in getModuleContent(module.type)" :key="index" class="item">
                <div class="item-header">
                  <span class="item-title">{{ item.company || '公司名称' }}</span>
                  <span class="item-date">
                    {{ item.start_date || '开始' }} - {{ item.end_date || '结束' }}
                  </span>
                </div>
                <div class="item-subtitle">{{ item.position || '职位' }}</div>
                <p v-if="item.description" class="item-description">{{ item.description }}</p>
                <ul v-if="item.highlights && item.highlights.length" class="highlights">
                  <li v-for="(hl, hlIndex) in item.highlights" :key="hlIndex" v-if="hl">
                    {{ hl }}
                  </li>
                </ul>
              </div>
              <el-empty v-if="getModuleContent(module.type).length === 0" description="暂无工作经历" :image-size="60" />
            </div>

            <div v-else-if="module.type === 'skills'" class="module">
              <h2 class="module-title" :style="{ borderColor: primaryColor }">
                <span :style="{ color: primaryColor }">{{ module.name }}</span>
              </h2>
              <div v-for="(item, index) in getModuleContent(module.type)" :key="index" class="item">
                <div v-if="item.category || item.items" class="skill-category">
                  <strong v-if="item.category">{{ item.category }}：</strong>
                  <span>{{ item.items || '' }}</span>
                </div>
              </div>
              <div v-if="getModuleContent(module.type).length === 0" class="skill-category">
                <span>暂无技能信息</span>
              </div>
            </div>

            <div v-else-if="module.type === 'projects'" class="module">
              <h2 class="module-title" :style="{ borderColor: primaryColor }">
                <span :style="{ color: primaryColor }">{{ module.name }}</span>
              </h2>
              <div v-for="(item, index) in getModuleContent(module.type)" :key="index" class="item">
                <div class="item-header">
                  <span class="item-title">{{ item.name || '项目名称' }}</span>
                  <span class="item-date">
                    {{ item.start_date || '开始' }} - {{ item.end_date || '结束' }}
                  </span>
                </div>
                <div class="item-subtitle">{{ item.role || '担任角色' }}</div>
                <p v-if="item.description" class="item-description">{{ item.description }}</p>
                <ul v-if="item.highlights && item.highlights.length" class="highlights">
                  <li v-for="(hl, hlIndex) in item.highlights" :key="hlIndex" v-if="hl">
                    {{ hl }}
                  </li>
                </ul>
              </div>
              <el-empty v-if="getModuleContent(module.type).length === 0" description="暂无项目经历" :image-size="60" />
            </div>

            <div v-else-if="module.type === 'summary' || module.type === 'evaluation'" class="module">
              <h2 class="module-title" :style="{ borderColor: primaryColor }">
                <span :style="{ color: primaryColor }">{{ module.name }}</span>
              </h2>
              <p class="text-content">{{ getModuleContent(module.type) || '暂无内容，请在左侧编辑' }}</p>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useResumeStore } from '@/stores/resume'

const store = useResumeStore()

const visibleModules = computed(() => store.visibleModules)
const styleConfig = computed(() => store.styleConfig)
const content = computed(() => store.content)

const primaryColor = computed(() => styleConfig.value.primary_color)

const containerStyle = computed(() => ({
  '--primary-color': styleConfig.value.primary_color,
  '--font-size': `${styleConfig.value.font_size}px`,
  '--line-height': styleConfig.value.line_height,
  '--page-margin': `${styleConfig.value.page_margin}px`
}))

const contentStyle = computed(() => ({
  fontSize: `${styleConfig.value.font_size}px`,
  lineHeight: styleConfig.value.line_height,
  padding: `${styleConfig.value.page_margin}px`
}))

const basicInfo = computed(() => content.value.basic || {})

function getModuleContent(type) {
  return content.value[type] || (type === 'summary' || type === 'evaluation' ? '' : [])
}
</script>

<style lang="scss" scoped>
.resume-preview {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f0f2f5;

  .preview-header {
    padding: 16px;
    background: #fff;
    border-bottom: 1px solid #e4e7ed;
    font-weight: 600;
    font-size: 14px;
    flex-shrink: 0;
  }

  .preview-body {
    flex: 1;
    overflow: auto;
    padding: 20px;
    display: flex;
    justify-content: center;

    .a4-container {
      width: 210mm;
      min-height: 297mm;
      background: #fff;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
      transform-origin: top center;
      transform: scale(0.65);
      margin-bottom: -35%;

      @media print {
        transform: scale(1);
        box-shadow: none;
        width: 100%;
        min-height: auto;
        margin-bottom: 0;
      }
    }

    .resume-content {
      color: #303133;
      font-size: var(--font-size);
      line-height: var(--line-height);

      .module {
        margin-bottom: calc(var(--margin-top, 20) * 0.8px);

        &.basic-module {
          text-align: center;
          padding-bottom: 20px;
          border-bottom: 2px solid var(--primary-color, #409eff);
          margin-bottom: calc(var(--margin-top, 20) * 0.8px);

          .name {
            font-size: calc(var(--font-size) * 2);
            font-weight: bold;
            margin-bottom: 12px;
            line-height: 1.3;
          }

          .basic-info-row {
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-bottom: 8px;
            font-size: calc(var(--font-size) - 1px);
            color: #606266;
            flex-wrap: wrap;
            line-height: var(--line-height);

            .info-item {
              display: inline-flex;
              align-items: center;
            }

            .link {
              color: var(--primary-color, #409eff);
              text-decoration: none;

              &:hover {
                text-decoration: underline;
              }
            }
          }
        }

        .module-title {
          font-size: calc(var(--font-size) * 1.3);
          font-weight: bold;
          padding-bottom: 8px;
          margin-bottom: 12px;
          border-bottom: 2px solid;
          line-height: 1.4;
        }

        .item {
          margin-bottom: 16px;
          line-height: var(--line-height);

          .item-header {
            display: flex;
            justify-content: space-between;
            align-items: baseline;
            margin-bottom: 4px;
            flex-wrap: wrap;
            gap: 8px;

            .item-title {
              font-weight: bold;
              font-size: calc(var(--font-size) + 1px);
            }

            .item-date {
              color: #909399;
              font-size: calc(var(--font-size) - 1px);
            }
          }

          .item-subtitle {
            color: #606266;
            margin-bottom: 8px;
            font-size: var(--font-size);
          }

          .item-description {
            color: #606266;
            margin-bottom: 8px;
            white-space: pre-wrap;
            font-size: var(--font-size);
            line-height: var(--line-height);
          }

          .highlights {
            margin: 8px 0;
            padding-left: 20px;
            font-size: var(--font-size);

            li {
              color: #606266;
              margin-bottom: 4px;
              line-height: var(--line-height);
            }
          }

          .skill-category {
            margin-bottom: 8px;
            color: #606266;
            font-size: var(--font-size);
            line-height: var(--line-height);

            strong {
              color: #303133;
            }
          }
        }

        .text-content {
          color: #606266;
          white-space: pre-wrap;
          line-height: var(--line-height);
          font-size: var(--font-size);
        }

        :deep(.el-empty) {
          padding: 20px 0;

          .el-empty__description {
            color: #c0c4cc;
            font-size: calc(var(--font-size) - 1px);
          }
        }
      }
    }
  }
}

@media print {
  body * {
    visibility: hidden;
  }

  .resume-preview,
  .resume-preview * {
    visibility: visible;
  }

  .resume-preview {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    background: #fff;
  }

  .preview-header {
    display: none !important;
  }

  .preview-body {
    padding: 0 !important;
    overflow: visible !important;
  }
}
</style>
