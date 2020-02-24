/*

Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package commands

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/net/context"

	"github.com/zdnscloud/lvmd/parser"
)

// ListLV lists lvm volumes
func ListLV(ctx context.Context, listspec string) ([]*parser.LV, error) {
	cmd := exec.Command("lvs", "--units=b", "--separator=<:SEP:>", "--nosuffix", "--noheadings",
		"-o", "lv_name,lv_size,lv_uuid,lv_attr,copy_percent,lv_kernel_major,lv_kernel_minor,lv_tags", "--nameprefixes", "-a", listspec)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	outStr := strings.TrimSpace(string(out))
	outLines := strings.Split(outStr, "\n")
	lvs := make([]*parser.LV, len(outLines))
	for i, line := range outLines {
		line = strings.TrimSpace(line)
		lv, err := parser.ParseLV(line)
		if err != nil {
			return nil, err
		}
		lvs[i] = lv
	}
	return lvs, nil
}

func CreateThinPoolUseAllSize(ctx context.Context, vg string, pool string) (string, error) {
	args := []string{"-v", "-l", "100%FREE", "--thinpool", pool, vg, "-y"}
	cmd := exec.Command("lvcreate", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CreateLV creates a new volume
func CreateThinLV(ctx context.Context, vg string, name string, size uint64, mirrors uint32, tags []string) (string, error) {
	if size == 0 {
		return "", errors.New("size must be greater than 0")
	}

	args := []string{"--thin", "-v", "-n", name, "-V", fmt.Sprintf("%db", size)}
	if mirrors > 0 {
		args = append(args, "-m", fmt.Sprintf("%d", mirrors), "--nosync")
	}
	for _, tag := range tags {
		args = append(args, "--add-tag", tag)
	}
	args = append(args, vg)
	cmd := exec.Command("lvcreate", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CreateLV creates a new volume
func CreateLV(ctx context.Context, vg string, name string, size uint64, mirrors uint32, tags []string) (string, error) {
	if size == 0 {
		return "", errors.New("size must be greater than 0")
	}

	args := []string{"-v", "-n", name, "-L", fmt.Sprintf("%db", size)}
	if mirrors > 0 {
		args = append(args, "-m", fmt.Sprintf("%d", mirrors), "--nosync")
	}
	for _, tag := range tags {
		args = append(args, "--add-tag", tag)
	}
	args = append(args, vg)
	cmd := exec.Command("lvcreate", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ProtectedTagName is a tag that prevents RemoveLV & RemoveVG from removing a volume
const ProtectedTagName = "protected"

// RemoveLV removes a volume
func RemoveLV(ctx context.Context, vg string, name string) (string, error) {
	lvs, err := ListLV(ctx, fmt.Sprintf("%s/%s", vg, name))
	if err != nil {
		return "", fmt.Errorf("failed to list LVs: %v", err)
	}
	if len(lvs) != 1 {
		return "", fmt.Errorf("expected 1 LV, got %d", len(lvs))
	}
	for _, tag := range lvs[0].Tags {
		if tag == ProtectedTagName {
			return "", errors.New("volume is protected")
		}
	}

	cmd := exec.Command("lvremove", "-v", "-f", fmt.Sprintf("%s/%s", vg, name))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CloneLV clones a volume via dd
func CloneLV(ctx context.Context, src, dest string) (string, error) {
	// FIXME(farcaller): bloody insecure. And broken.
	cmd := exec.Command("dd", fmt.Sprintf("if=%s", src), fmt.Sprintf("of=%s", dest), "bs=4M")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ResizeLV(ctx context.Context, vg string, name string, size uint64) (string, error) {
	cmd := exec.Command("lvresize", "-L", fmt.Sprintf("%db", size), "-v", fmt.Sprintf("%s/%s", vg, name))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ResizeLVe2fsck(ctx context.Context, vg string, name string) (string, error) {
	cmd := exec.Command("e2fsck", "-f", "-y", fmt.Sprintf("/dev/%s/%s", vg, name))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ResizeLV2fs(ctx context.Context, vg string, name string) (string, error) {
	cmd := exec.Command("resize2fs", fmt.Sprintf("/dev/%s/%s", vg, name))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ListVG(ctx context.Context) ([]*parser.VG, error) {

	cmd := exec.Command("vgs", "--units=b", "--separator=<:SEP:>", "--nosuffix", "--noheadings",
		"-o", "vg_name,vg_size,vg_free,vg_uuid,vg_tags", "--nameprefixes", "-a")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	outStr := strings.TrimSpace(string(out))
	outLines := strings.Split(outStr, "\n")
	vgs := make([]*parser.VG, len(outLines))
	for i, line := range outLines {
		line = strings.TrimSpace(line)
		vg, err := parser.ParseVG(line)
		if err != nil {
			return nil, err
		}
		vgs[i] = vg
	}
	return vgs, nil
}

func ExtendVG(ctx context.Context, name string, physicalVolume string) (string, error) {
	cmd := exec.Command("vgextend", name, physicalVolume)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ReduceVG(ctx context.Context, name string, physicalVolume string) (string, error) {
	cmd := exec.Command("vgreduce", name, physicalVolume)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func CreateVG(ctx context.Context, name string, physicalVolume string, tags []string) (string, error) {
	args := []string{name, physicalVolume, "-v"}
	for _, tag := range tags {
		args = append(args, "--add-tag", tag)
	}
	cmd := exec.Command("vgcreate", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func RemoveVG(ctx context.Context, name string) (string, error) {
	vgs, err := ListVG(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list VGs: %v", err)
	}
	var vg *parser.VG
	for _, v := range vgs {
		if v.Name == name {
			vg = v
			break
		}
	}
	if vg == nil {
		return "", fmt.Errorf("could not find vg to delete")
	}
	for _, tag := range vg.Tags {
		if tag == ProtectedTagName {
			return "", errors.New("volume is protected")
		}
	}

	cmd := exec.Command("vgremove", "-v", "-f", name)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func AddTagLV(ctx context.Context, vg string, name string, tags []string) (string, error) {
	lvs, err := ListLV(ctx, fmt.Sprintf("%s/%s", vg, name))
	if err != nil {
		return "", fmt.Errorf("failed to list LVs: %v", err)
	}
	if len(lvs) != 1 {
		return "", fmt.Errorf("expected 1 LV, got %d", len(lvs))
	}

	args := make([]string, 0)

	for _, tag := range tags {
		args = append(args, "--addtag", tag)
	}

	args = append(args, fmt.Sprintf("%s/%s", vg, name))

	cmd := exec.Command("lvchange", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func RemoveTagLV(ctx context.Context, vg string, name string, tags []string) (string, error) {
	lvs, err := ListLV(ctx, fmt.Sprintf("%s/%s", vg, name))
	if err != nil {
		return "", fmt.Errorf("failed to list LVs: %v", err)
	}
	if len(lvs) != 1 {
		return "", fmt.Errorf("expected 1 LV, got %d", len(lvs))
	}

	args := make([]string, 0)

	for _, tag := range tags {
		args = append(args, "--deltag", tag)
	}

	args = append(args, fmt.Sprintf("%s/%s", vg, name))

	cmd := exec.Command("lvchange", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func CreatePV(ctx context.Context, block string) (string, error) {
	args := []string{block, "-y", "-v"}
	cmd := exec.Command("pvcreate", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func RemovePV(ctx context.Context, block string) (string, error) {
	args := []string{block, "-y", "-v"}
	cmd := exec.Command("pvremove", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ListPV(ctx context.Context) ([]*parser.PV, error) {
	cmd := exec.Command("pvs", "--units=b", "--separator=<:SEP:>", "--nosuffix", "--noheadings",
		"-o", "pv_name,pv_size,pv_used,pv_free,pv_fmt,pv_uuid", "--nameprefixes", "-a")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	outStr := strings.TrimSpace(string(out))
	outLines := strings.Split(outStr, "\n")
	pvs := make([]*parser.PV, len(outLines))
	for i, line := range outLines {
		line = strings.TrimSpace(line)
		pv, err := parser.ParsePV(line)
		if err != nil {
			return nil, err
		}
		pvs[i] = pv
	}
	return pvs, nil
}

func Validate(ctx context.Context, block string) (bool, error) {
	cmd := exec.Command("udevadm", "info", "--query=property", block)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	out1Str := strings.TrimSpace(string(out))
	if strings.Contains(out1Str, "ID_PART_TABLE") || strings.Contains(out1Str, "ID_FS_TYPE") {
		return false, nil
	}

	out, err = exec.Command("blkid").Output()
	if err != nil {
		return false, err
	}
	outputs := strings.Split(string(out), "\n")
	for _, l := range outputs {
		if !strings.Contains(l, block) {
			continue
		}
		if strings.Contains(l, "PTTYPE=") || strings.Contains(l, "TYPE=") {
			return false, nil
		}
	}
	return true, nil
}

func Destory(ctx context.Context, block string) (string, error) {
	cmd := exec.Command("wipefs", "-af", block)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func Match(ctx context.Context, block string) string {
	cmd := exec.Command("pvs", "--noheadings", "--separator=#", "--nosuffix", block)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	outStr := strings.TrimSpace(string(out))
	return strings.Split(outStr, "#")[1]
}

func GetPVNum(ctx context.Context, name string) (string, error) {
	cmd := exec.Command("vgs", "--noheadings", "--separator=#", "--nosuffix", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "0", err
	}
	outStr := strings.TrimSpace(string(out))
	return strings.Split(outStr, "#")[1], nil
}
